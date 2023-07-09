package rest

import (
	"net/http"

	"github.com/AnthonyHewins/aq/internal/d"
	"github.com/AnthonyHewins/aq/internal/d/census"
	"github.com/AnthonyHewins/aq/internal/middleware"
	"github.com/AnthonyHewins/aq/internal/svc"
	"github.com/AnthonyHewins/aq/internal/svc/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type server struct {
	tracer trace.Tracer
	logger log.Logger

	reader *sqlx.DB
	writer *sqlx.DB

	insightClient *d.Client
	censusClient  *census.DBClient
	authClient    *auth.Client
	controller    *svc.Controller
}

func extractClaims(c *gin.Context) {
	bearer := c.Request.Header.Get("authorization") // AWS lowercases
	if bearer == "" {
		bearer = c.Request.Header.Get("Authorization") // localhost may not
	}

	claims, err := middleware.JWTBearerValidate(bearer)
	if err != nil {
		baseErrHandler(c, err)
		return
	}

	c.Set("x-sub", claims.Subject)
	c.Next()
}

func NewServer(logger log.Logger, httpClient *http.Client, dbReader, dbWriter *sqlx.DB, noMiddlewares bool) *gin.Engine {
	ginServer := gin.Default()

	if !noMiddlewares {
		ginServer.Use(extractClaims)
	}

	s := &server{
		tracer:        otel.Tracer("rest"),
		logger:        logger,
		reader:        dbReader,
		writer:        dbWriter,
		insightClient: d.NewClient("rest-insight-svc", logger, httpClient, dbReader),
		censusClient:  census.NewDBClient("rest-jurisdiction-svc", logger, dbReader),
		controller:    svc.NewController("controller", logger, dbReader, dbWriter),
		authClient:    auth.NewClient("rest-auth-client", logger, dbReader, dbWriter),
	}

	v1 := ginServer.Group("/api/v1")

	{ // /api/mobile-homes
		mobileHome := v1.Group("/mobile-homes")
		mobileHome.GET("", s.mobileHomeIndex())

		rest := mobileHome.Group("/:id")
		rest.GET("", s.mobileHomeGet())
	}

	{
		//		residential := v1.Group("/residential")
		//		residential.GET("", s.residentialIndex())
		//
		//		rest := residential.Group("/:id")
		//		rest.GET("", s.residentialGet())
	}

	{ // /api/insights
		insights := v1.Group("/insights")
		insights.POST("/location", s.location())
	}

	{ // /api/jurisdictions
		jurisdictions := v1.Group("/jurisdictions")
		jurisdictions.GET("", s.indexJurisdictions())
		jurisdictions.POST("", s.purchaseJurisdiction())
	}

	return ginServer
}
