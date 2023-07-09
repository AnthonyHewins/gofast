package grpcserver

import (
	"context"
	"net/http"
	"time"

	"github.com/AnthonyHewins/aq/gen/go/api/proto/insightsvc/v1"
	"github.com/AnthonyHewins/aq/gen/go/api/proto/jurisdictionsvc/v1"
	"github.com/AnthonyHewins/aq/gen/go/api/proto/mobilehomesvc/v1"
	"github.com/AnthonyHewins/aq/internal/api"
	"github.com/AnthonyHewins/aq/internal/d"
	"github.com/AnthonyHewins/aq/internal/middleware"
	"github.com/AnthonyHewins/aq/internal/svc"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type ctxKey string

const (
	_userID ctxKey = "id"
)

var (
	errNotAdmin   = status.Error(codes.PermissionDenied, "this route is only available to admins")
	errMissingReq = status.Error(codes.InvalidArgument, "missing request")
)

var baseUnaryMiddlewares = []grpc.UnaryServerInterceptor{
	authorizationMiddleware, // strips metadata needed for request
}

func fetchKey(m metadata.MD, key string) (string, error) {
	switch vals := m.Get(key); len(vals) {
	case 0:
		return "", nil
	case 1:
		return vals[0], nil
	default:
		return "", status.Errorf(codes.InvalidArgument, "metadata value '%s' should only have one value, but got %d: %v", key, len(vals), vals)
	}
}

func authorizationMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing gRPC metadata in request")
	}

	bearer, err := fetchKey(m, "authorization")
	if err != nil {
		return nil, err
	}

	claims, err := middleware.JWTBearerValidate(bearer)
	if err != nil {
		apiErr, ok := err.(api.APIErr)
		if !ok {
			return nil, err
		}
		return nil, status.Error(apiErr.Code(), apiErr.Error())
	}

	if claims.Subject == uuid.Nil {
		return nil, status.Error(codes.InvalidArgument, "invalid UUID passed for user ID: got nil")
	}

	return handler(context.WithValue(ctx, _userID, claims.Subject), req)
}

type server struct {
	logger log.Logger
	tracer trace.Tracer

	reader *sqlx.DB
	writer *sqlx.DB

	controller    *svc.Controller
	insightClient *d.Client
}

// NewServer creates a new server. Pass an empty string to traceName to not add any trace middleware
func NewServer(traceName string, l log.Logger, reader, writer *sqlx.DB, httpClient *http.Client) *grpc.Server {
	if reader == nil || writer == nil {
		level.Error(l).Log(
			"msg", "fatal error: dependency missing for grpc server. Server cannot start",
			"reader", reader,
			"writer", writer,
		)
		return nil
	}

	s := &server{
		logger:        l,
		tracer:        otel.Tracer(traceName),
		reader:        reader,
		writer:        writer,
		controller:    svc.NewController(traceName+" svc controller", l, reader, writer),
		insightClient: d.NewClient(traceName+" insight client", l, httpClient, reader),
	}

	// only if tracing was specified should you add this
	if traceName != "" {
		baseUnaryMiddlewares = append(baseUnaryMiddlewares, otelgrpc.UnaryServerInterceptor())
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Second,
			MaxConnectionAge:  5 * time.Minute,
			Timeout:           20 * time.Second,
		}),

		grpc.ChainUnaryInterceptor(baseUnaryMiddlewares...),

		// Stream requests
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	insightsvc.RegisterInsightServiceServer(grpcServer, s)
	mobilehomesvc.RegisterMobileHomeServiceServer(grpcServer, s)
	jurisdictionsvc.RegisterJurisdictionServiceServer(grpcServer, s)

	reflection.Register(grpcServer)
	return grpcServer
}
