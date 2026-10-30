FROM golang:1.26-alpine AS build
RUN apk add --no-cache git protobuf protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6 \
 && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
WORKDIR /src
COPY services/gateway/go.mod services/gateway/go.sum ./
RUN go mod download
COPY proto /proto
COPY services/gateway/ ./
RUN protoc -I /proto \
      --go_out=. --go_opt=module=github.com/connect/gateway \
      --go-grpc_out=. --go-grpc_opt=module=github.com/connect/gateway \
      /proto/connect/v1/gateway.proto \
 && CGO_ENABLED=0 go build -o /gateway ./cmd/gateway

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /gateway /gateway
EXPOSE 8080 50051
ENTRYPOINT ["/gateway"]
