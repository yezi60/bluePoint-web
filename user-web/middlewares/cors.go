package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 用于跨域问题下处理options请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Token,x-token,Authorization")
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE,PATCH,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control,Access-Control-Allow-Headers,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

	}
}
