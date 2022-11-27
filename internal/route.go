package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"year-end/internal/handler"
	"year-end/internal/middleware"
)

func generateRoute() *gin.Engine {
	gin.SetMode(viper.GetString("gin.mode"))
	var engine = gin.Default()
	engine.POST("/api/v1/muxi/year-task/login", handler.LoginAPI)
	
	auth := engine.Group("/api/v1/muxi/year-task/user")
	{
		auth.Use(middleware.Parse())
		// 一个用户的数据统计
		auth.GET("/self/info", handler.GetData)
		// 全校学生的综合统计
		auth.GET("/all/info")
	}
	
	return engine
}

func StartHTTP() *gin.Engine {
	engine := generateRoute()
	engine.Run(":" + viper.GetString("gin.port"))
	return engine
}
