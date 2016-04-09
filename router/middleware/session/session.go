package session

import (
	"errors"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/robzienert/gin-middleware/oauth"
	"github.com/robzienert/lever/shared/strutil"
)

// User is a convenience net.Context function for retrieving the authorized
// OAuth user.
func User(c *gin.Context) *oauth.User {
	token := oauth.Token(c)
	if token == nil || token.User == nil {
		return nil
	}
	return token.User
}

// AuditActor is a convenience net.Context function that wraps User for auditing
// purposes.
func AuditActor(c *gin.Context) string {
	token := oauth.Token(c)
	if token == nil {
		logrus.Error("no oauth token to determine audit actor: this should never happen")
		return "unknown"
	}
	if token.User != nil {
		return token.User.Username
	}
	return token.ClientID
}

// AuthFeatureState contains the logic to determine if a request to process
// feature state should be allowed. If the request is allowed, this function
// will return a nil error.
//
// 1. Service-scoped requests are allowed-all.
// 2. Mobile-scoped requests must have a user object and if "Actors" have been
// passed in the request, the user's username may be the only actor value.
func AuthFeatureState() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := oauth.Token(c)
		if token == nil {
			c.Error(errors.New("no oauth token to determine feature state authorization")).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if strutil.StringInSlice("service", token.Scopes) {
			c.Next()
			return
		}
		user := User(c)
		if user == nil {
			c.Error(errors.New("incorrect scopes and no user found in oauth token")).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		actors := c.Query("actors")
		if actors == "" || user.Username == actors {
			c.Next()
			return
		}
		c.Error(errors.New("user cannot request other actor feature states")).SetType(gin.ErrorTypePublic)
		c.AbortWithStatus(http.StatusForbidden)
	}
}
