package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robzienert/http-healthcheck"
)

// GetHealthStatus returns the health status of the service.
func GetHealthStatus(c *gin.Context) {
	status := healthcheck.FromContext(c).Status()
	resp := healthcheck.MarshalHealthStatusResponse(status)
	if status.Healthy {
		c.IndentedJSON(http.StatusOK, resp)
	} else {
		c.IndentedJSON(http.StatusInternalServerError, resp)
	}
}
