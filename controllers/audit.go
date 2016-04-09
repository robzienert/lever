package controllers

import (
	"net/http"

	"github.com/robzienert/lever/api"
	"github.com/robzienert/lever/model"
	"github.com/robzienert/lever/store"
	"github.com/gin-gonic/gin"
)

// GetAuditIndex returns a list of audit breadcrumbs.
func GetAuditIndex(c *gin.Context) {
	breadcrumbs, err := store.GetBreadcrumbList(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if breadcrumbs == nil {
		breadcrumbs = make([]*model.Breadcrumb, 0)
	}
	c.IndentedJSON(http.StatusOK, api.GetAuditResponse{
		Breadcrumbs: breadcrumbs,
	})
}
