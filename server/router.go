package server

import (
	"embed"
	"io/fs"
	"moeCounter/database"
	"moeCounter/server/controller"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 注册路由到现有路由器
func RegisterRoutes(port int, dbFile string, publicFS embed.FS) {
	router := gin.Default()
	router.Use(gin.Recovery())

	// 初始化数据库
	if err := database.InitDB(dbFile); err != nil {
		panic("数据库初始化失败: " + err.Error())
	}

	// 从嵌入文件系统获取public子目录
	fsPublic, err := fs.Sub(publicFS, "public")
	if err != nil {
		panic("无法获取public子目录: " + err.Error())
	}

	// 获取assets子目录
	fsAssets, err := fs.Sub(fsPublic, "assets")
	if err != nil {
		panic("无法获取assets子目录: " + err.Error())
	}

	// 注册根路径路由，返回首页
	router.GET("/", func(c *gin.Context) {
		data, err := fs.ReadFile(fsPublic, "index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// 注册静态文件服务
	router.StaticFS("/assets", http.FS(fsAssets))

	// 计数器接口
	router.GET("/counter", controller.CounterHandler(fsPublic))

	// 添加主题列表接口
	router.GET("/themes", controller.ThemeListHandler(fsPublic))

	// 使用传入的端口
	if err := router.Run(":" + strconv.Itoa(port)); err != nil {
		router.Run() // 使用随机端口
	}
}
