package middleware

import (
	"go-gin-hexagonal/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

func CSRFMiddleware() gin.HandlerFunc {
	csrfKey := utils.GeneratePassword(32, false)
	csrfProtection := csrf.Protect([]byte(csrfKey), csrf.Secure(false), csrf.MaxAge(600), csrf.CookieName("CSRF-TOKEN"))

	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		if strings.HasPrefix(c.Request.URL.Path, "/api/v1/auth") {
			c.Next()
			return
		}

		csrfProtection(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		}))
	}
}
