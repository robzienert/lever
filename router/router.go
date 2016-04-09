package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/expvar"
	"github.com/gin-gonic/gin"
	gomhttp "github.com/rakyll/gom/http"
	"github.com/robzienert/gin-middleware/correlationid"
	"github.com/robzienert/gin-middleware/err"
	"github.com/robzienert/gin-middleware/header"
	"github.com/robzienert/gin-middleware/localhost"
	"github.com/robzienert/gin-middleware/oauth"
	"github.com/robzienert/lever/controllers"
	"github.com/robzienert/lever/router/middleware/session"
)

var (
	consumerScope = []string{"service", "mobile"}
	serviceScope  = []string{"service"}
)

// Load will setup the HTTP engine, middleware and router.
func Load(oauthValidator oauth.TokenValidator, middleware ...gin.HandlerFunc) http.Handler {
	e := gin.Default()

	e.Use(correlationid.SetRequestUUID(correlationid.CorrelationHeader))
	e.Use(err.ErrorHandler())
	e.Use(correlationid.RequestLogger(time.RFC3339, true))
	e.Use(header.NoCache())
	e.Use(header.Secure())
	e.Use(middleware...)

	api := e.Group("/api")
	{
		api.Use(oauth.BearerTokenAuth(oauthValidator))
		mustConsumer := oauth.MustScope(consumerScope)
		mustService := oauth.MustScope(serviceScope)
		authFeatureState := session.AuthFeatureState()

		features := api.Group("/features")
		{
			// There's a bug in gin that disallows us from having a /state endpoint
			// as well as a wildcard /:key. For that reason, I moved the batch
			// endpoint to POST /. It's not great, but a lot simpler than doing a
			// workaround.
			features.GET("", mustConsumer, controllers.GetAllFeatures)
			features.POST("", mustConsumer, authFeatureState, controllers.PostBatchFeatureState)
			features.GET("/:key", mustConsumer, controllers.GetFeature)
			features.PUT("/:key", mustService, controllers.PutFeature)
			features.DELETE("/:key", mustService, controllers.DeleteFeature)
			features.GET("/:key/state", mustConsumer, authFeatureState, controllers.GetFeatureState)
		}
		api.GET("/audit", mustService, controllers.GetAuditIndex)
	}

	e.GET("/status", controllers.GetHealthStatus)

	e.GET("/debug/vars", localhost.MustLocal(), expvar.Handler())
	if gin.IsDebugging() {
		e.GET("/debug/_gom", gin.WrapF(gomhttp.Handler()))
	}

	return e
}
