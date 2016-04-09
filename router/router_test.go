package router

import (
	"fmt"
	"testing"

	"github.com/robzienert/lever/controllers"
	"github.com/robzienert/lever/shared/httputil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func assertRouteExists(t *testing.T, routes gin.RoutesInfo, expectedMethod string, expectedPath string, expectedHandler interface{}) {
	assert.True(t, httputil.RouteExists(routes, expectedMethod, expectedPath, expectedHandler), fmt.Sprintf("%s %s route does not exist", expectedMethod, expectedPath))
}

type RouterTestSuite struct {
	suite.Suite
}

func (suite *RouterTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (suite *RouterTestSuite) TestLoad_Routes() {
	router, ok := Load(nil).(*gin.Engine)
	assert.True(suite.T(), ok, "could not cast http.Handler has *gin.Engine")

	routes := router.Routes()
	assertRouteExists(suite.T(), routes, "GET", "/api/audit", controllers.GetAuditIndex)
	assertRouteExists(suite.T(), routes, "GET", "/api/features", controllers.GetAllFeatures)
	assertRouteExists(suite.T(), routes, "POST", "/api/features", controllers.PostBatchFeatureState)
	assertRouteExists(suite.T(), routes, "GET", "/api/features/:key", controllers.GetFeature)
	assertRouteExists(suite.T(), routes, "PUT", "/api/features/:key", controllers.PutFeature)
	assertRouteExists(suite.T(), routes, "DELETE", "/api/features/:key", controllers.DeleteFeature)
	assertRouteExists(suite.T(), routes, "GET", "/api/features/:key/state", controllers.GetFeatureState)
	assertRouteExists(suite.T(), routes, "GET", "/status", controllers.GetHealthStatus)
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}
