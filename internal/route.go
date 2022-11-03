package internal

import "github.com/gin-gonic/gin"

func generateRoute() *gin.Engine {
	var engine = gin.Default()
	
	return engine
}

func StartHTTP(port string) *gin.Engine {
	engine := generateRoute()
	engine.Run(":" + port)
	return engine
}
