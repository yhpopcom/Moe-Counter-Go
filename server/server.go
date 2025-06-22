package server

import (
	"embed"
	"moeCounter/internal/database"

	"github.com/gin-gonic/gin"
)

// 初始化路由
func InitRouter(port int, dbFile string, publicFS embed.FS, debug bool) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	if debug {
		router.Use(gin.Logger())
	}

	// 初始化数据库
	if err := database.InitDB(dbFile, debug); err != nil {
		panic("数据库初始化失败: " + err.Error())
	}

	// 注册基础路由
	registerBaseRoutes(router, publicFS)

	// 注册API路由组
	apiGroup := router.Group("/api")
	registerAPIRoutes(apiGroup)

	return router
}
