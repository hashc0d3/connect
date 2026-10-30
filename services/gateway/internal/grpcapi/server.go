package grpcapi

import (
	"context"
	"errors"
	"net"

	connectv1 "github.com/connect/gateway/gen/connect/v1"
	"github.com/connect/gateway/internal/authn"
	"github.com/connect/gateway/internal/sms"
	"github.com/connect/gateway/internal/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type userKey struct{}

type Server struct {
	connectv1.UnimplementedSmsGatewayServer
	Store *store.Store
	SMS   *sms.Service
}

func Listen(addr string, db *store.Store, service *sms.Service) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}
	srv := &Server{Store: db, SMS: service}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(srv.authenticate))
	connectv1.RegisterSmsGatewayServer(grpcServer, srv)
	reflection.Register(grpcServer)
	return grpcServer, lis, nil
}

func (s *Server) authenticate(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	header := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("authorization")
		if len(values) > 0 {
			header = values[0]
		}
	}
	user, err := authn.ByAuthorization(ctx, s.Store, header)
	if err != nil || user == nil {
		if err != nil && !errors.Is(err, authn.ErrUnauthenticated) {
			return nil, status.Error(codes.Internal, "Внутренняя ошибка")
		}
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated.")
	}
	return handler(context.WithValue(ctx, userKey{}, user.ID), req)
}

func userID(ctx context.Context) (string, error) {
	id, _ := ctx.Value(userKey{}).(string)
	if id == "" {
		return "", status.Error(codes.Unauthenticated, "Unauthenticated.")
	}
	return id, nil
}
