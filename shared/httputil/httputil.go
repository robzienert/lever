package httputil

import (
	"github.com/robzienert/lever/shared/reflectutil"
	"github.com/gin-gonic/gin"
)

// RouteExists is a utility function for determining if a gin route exists at
// the given path and method, with the expected handler.
func RouteExists(routes gin.RoutesInfo, expectedMethod string, expectedPath string, expectedHandler interface{}) bool {
	for _, r := range routes {
		if r.Method == expectedMethod && r.Path == expectedPath && r.Handler == reflectutil.FuncName(expectedHandler) {
			return true
		}
	}
	return false
}
