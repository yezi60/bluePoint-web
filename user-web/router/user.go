package router

import (
	"bluepoint_api/user-web/api"
	"bluepoint_api/user-web/middlewares"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

//
func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("配置用户相关的url")
	{
		UserRouter.GET("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
		UserRouter.POST("register", api.Register)
		UserRouter.GET("detail", middlewares.JWTAuth(), api.GetUserDetail)
		UserRouter.PATCH("update", middlewares.JWTAuth(), api.UpdateUser)
	}
}
