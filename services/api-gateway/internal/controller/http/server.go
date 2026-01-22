package http

import (
	"api_gateway/internal/config"
	"api_gateway/pkg/metrics"
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authv1 "github.com/RomanKovalev007/tutors_platform/api/gen/go/auth"
	groupv1 "github.com/RomanKovalev007/tutors_platform/api/gen/go/group"
	tasksv1 "github.com/RomanKovalev007/tutors_platform/api/gen/go/task"
	userv1 "github.com/RomanKovalev007/tutors_platform/api/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient interface {
	ValidateToken(ctx context.Context, token string) (userID string, err error)
	Close() error
}

type Server struct {
	srv        *http.Server
	authClient AuthClient
	mux        *http.ServeMux
	gwMux      *runtime.ServeMux
	cfg        *config.Config
}

func NewServer(authClient AuthClient, cfg *config.Config) *Server {
	mux := http.NewServeMux()
	gwMux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
		runtime.WithOutgoingHeaderMatcher(customHeaderMatcher),
		runtime.WithMetadata(extractUserID),
	)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.HTTPPort),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &Server{
		srv:        srv,
		authClient: authClient,
		mux:        mux,
		gwMux:      gwMux,
		cfg:        cfg,
	}
}

func (s *Server) RegisterHandlers() error {
	ctx := context.Background()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, s.gwMux, s.cfg.AuthGRPC, opts); err != nil {
		return fmt.Errorf("failed to register auth gateway: %w", err)
	}

	if err := groupv1.RegisterGroupsServiceHandlerFromEndpoint(ctx, s.gwMux, s.cfg.GroupGRPC, opts); err != nil {
		return fmt.Errorf("failed to register group gateway: %w", err)
	}

	if err := userv1.RegisterUserServiceHandlerFromEndpoint(ctx, s.gwMux, s.cfg.UserGRPC, opts); err != nil {
		return fmt.Errorf("failed to register user gateway: %w", err)
	}

	if err := tasksv1.RegisterTaskServiceHandlerFromEndpoint(ctx, s.gwMux, s.cfg.TasksGRPC, opts); err != nil {
		return fmt.Errorf("failed to register tasks gateway: %w", err)
	}

	s.mux.Handle("/v1/", metrics.MetricsMiddleware(s.AuthMiddleware(s.gwMux)))
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.Handle("/metrics", metrics.Handler())

	s.srv.Handler = s.mux
	return nil
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
