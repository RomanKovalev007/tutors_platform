package http

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

func extractUserID(ctx context.Context, req *http.Request) metadata.MD {
	userID := req.Header.Get("X-User-Id")
	if userID != "" {
		md := metadata.Pairs("x-user-id", userID)
		return md
	}
	return nil
}

func customHeaderMatcher(key string) (string, bool) {
	switch key {
	case "X-User-Id":
		return "x-user-id", true
	case "Authorization":
		return "authorization", true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
