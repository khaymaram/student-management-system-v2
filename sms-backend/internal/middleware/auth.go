package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sms-backend/internal/services"
	"strings"
)

func Authenticate(auth *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "authentication required"})
			return
		}
		claims, err := auth.Parse(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "invalid or expired session"})
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("subjectID", claims.SubjectID)
		c.Next()
	}
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role == "admin" || strings.HasPrefix(c.Request.URL.Path, "/api/auth/") {
			c.Next()
			return
		}
		path := strings.TrimPrefix(c.Request.URL.Path, "/api")
		method := c.Request.Method
		subject := c.GetString("subjectID")
		allowed := false

		if role == "student" {
			studentPrefix := "/students/" + subject
			financePath := "/finances/" + subject
			allowed = method == http.MethodGet && (strings.HasPrefix(path, "/courses") || strings.HasPrefix(path, "/professors") || path == "/majors" || path == studentPrefix || path == financePath || strings.HasPrefix(path, studentPrefix+"/enrollments"))
			allowed = allowed || (method == http.MethodPut && path == financePath)
			allowed = allowed || (strings.HasPrefix(path, studentPrefix+"/enrollments") && (method == http.MethodPost || method == http.MethodDelete))
		} else if role == "professor" {
			allowed = method == http.MethodGet && (strings.HasPrefix(path, "/courses") || strings.HasPrefix(path, "/students") || path == "/professors/"+subject || path == "/enrollments")
			allowed = allowed || ((method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete) && strings.Contains(path, "/enrollments"))
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "error": "you do not have permission to perform this action"})
			return
		}
		c.Next()
	}
}
