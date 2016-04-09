package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/robzienert/http-healthcheck"
	"github.com/robzienert/lever/router/middleware/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type mockHealthProvider struct {
	healthy error
}

func (p *mockHealthProvider) Name() string     { return "mock" }
func (p *mockHealthProvider) Start() error     { return nil }
func (p *mockHealthProvider) IsHealthy() error { return p.healthy }
func (p *mockHealthProvider) Close() error     { return nil }

type StatusTestSuite struct {
	suite.Suite
}

func (suite *StatusTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (suite *StatusTestSuite) TestGetStatus() {
	router := gin.New()

	healthMonitor := healthcheck.New(healthcheck.DefaultSupervisor, &mockHealthProvider{healthy: nil})
	{
		defer healthMonitor.Close()
		healthMonitor.Start()
	}

	router.Use(context.SetHealthMonitor(healthMonitor))
	router.GET("/status", GetHealthStatus)

	req, _ := http.NewRequest("GET", "/status", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), 200, resp.Code)

	expected, _ := json.Marshal(healthcheck.HealthStatusResponse{
		Status: map[string]string{
			"mock": "OK",
		},
	})
	assert.JSONEq(suite.T(), string(expected), resp.Body.String())
}

func TestStatusTestSuite(t *testing.T) {
	suite.Run(t, new(StatusTestSuite))
}
