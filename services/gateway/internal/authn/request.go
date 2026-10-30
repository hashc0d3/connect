package authn

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/connect/gateway/internal/store"
)

var ErrUnauthenticated = errors.New("unauthenticated")

func FromRequest(r *http.Request, db *store.Store, internalToken string) (*store.User, error) {
	if token := strings.TrimSpace(r.Header.Get("X-Internal-Token")); token != "" {
		if subtle.ConstantTimeCompare([]byte(token), []byte(internalToken)) != 1 {
			return nil, ErrUnauthenticated
		}
		user, err := db.UserByID(r.Context(), strings.TrimSpace(r.Header.Get("X-User-Id")))
		if err != nil || user == nil {
			if err != nil {
				return nil, err
			}
			return nil, ErrUnauthenticated
		}
		return user, nil
	}
	if email, key, ok := r.BasicAuth(); ok {
		return byKey(r.Context(), db, key, email)
	}
	header := r.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return byKey(r.Context(), db, strings.TrimSpace(header[7:]), "")
	}
	if strings.HasPrefix(strings.ToLower(header), "basic ") {
		raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(header[6:]))
		if err != nil {
			return nil, ErrUnauthenticated
		}
		email, key, _ := strings.Cut(string(raw), ":")
		return byKey(r.Context(), db, key, email)
	}
	return nil, ErrUnauthenticated
}

func ByAuthorization(ctx context.Context, db *store.Store, header string) (*store.User, error) {
	header = strings.TrimSpace(header)
	if strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return byKey(ctx, db, strings.TrimSpace(header[7:]), "")
	}
	return nil, ErrUnauthenticated
}

func byKey(ctx context.Context, db *store.Store, key, email string) (*store.User, error) {
	if strings.TrimSpace(key) == "" {
		return nil, ErrUnauthenticated
	}
	user, _, err := db.UserByKeyHash(ctx, HashKey(key))
	if err != nil || user == nil {
		if err != nil {
			return nil, err
		}
		return nil, ErrUnauthenticated
	}
	if email != "" && !strings.EqualFold(user.Email, email) {
		return nil, ErrUnauthenticated
	}
	return user, nil
}
