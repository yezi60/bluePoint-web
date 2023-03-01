package initialize

import (
	"bluepoint_api/user-web/middlewares"
	"bluepoint_api/user-web/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

//初始化router

func Routers() *gin.Engine {
	Router := gin.Default()

	// 服务注册，健康
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	// 全局的跨域中间件
	Router.Use(middlewares.Cors())

	// 全局的group,可以添加
	//ApiGroup := Router.Group("/u/v1")

	//加入kong之后，使用/u作为service的区别
	ApiGroup := Router.Group("/v1")

	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)

	return Router
}
