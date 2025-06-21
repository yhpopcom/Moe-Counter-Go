package controller

import (
	"moeCounter/public"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 主题列表处理函数
func ThemeListHandler(c *gin.Context) {
	// 使用全局public文件系统
	publicFS := public.Public
	themes, err := ListThemes(publicFS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法获取主题列表",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"themes": themes,
	})
}
