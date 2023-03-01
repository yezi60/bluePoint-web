package middlewares

import (
	"bluepoint_api/user-web/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// IsAdminAuth 用于校验管理员权限的中间件
func IsAdminAuth() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityID != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}

}
