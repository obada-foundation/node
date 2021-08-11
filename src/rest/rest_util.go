package rest

import (
	"context"
	"net/http"
)

type util struct{}

func (u *util) ping(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))

	return nil
}
