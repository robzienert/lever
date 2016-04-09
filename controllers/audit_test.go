package controllers

import (
	"encoding/json"
	"testing"

	"github.com/robzienert/lever/api"
	"github.com/robzienert/lever/model"
	"github.com/robzienert/lever/router/middleware/context"
	"github.com/robzienert/lever/store/memory"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"net/http"
	"net/http/httptest"
)

type AuditTestSuite struct {
	suite.Suite
}

func (suite *AuditTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (suite *AuditTestSuite) TestAuditGetIndex_Empty() {
	router := gin.New()

	memStore := memory.Load()
	router.Use(context.SetStore(memStore))
	router.GET("/audit", GetAuditIndex)

	req, _ := http.NewRequest("GET", "/audit", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), 200, resp.Code)

	expected, _ := json.Marshal(api.GetAuditResponse{Breadcrumbs: []*model.Breadcrumb{}})
	assert.JSONEq(suite.T(), string(expected), resp.Body.String())
}

func (suite *AuditTestSuite) TestAuditGetIndex_NotEmpty() {
	router := gin.New()

	crumb := model.NewBreadcrumb("foo", "bar")

	memStore := memory.Load()
	memStore.Breadcrumbs().Create(crumb)
	router.Use(context.SetStore(memStore))
	router.GET("/audit", GetAuditIndex)

	req, _ := http.NewRequest("GET", "/audit", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), 200, resp.Code)

	expected, _ := json.Marshal(api.GetAuditResponse{
		Breadcrumbs: []*model.Breadcrumb{
			crumb,
		},
	})
	assert.JSONEq(suite.T(), string(expected), resp.Body.String())
}

func TestAuditTestSuite(t *testing.T) {
	suite.Run(t, new(AuditTestSuite))
}
