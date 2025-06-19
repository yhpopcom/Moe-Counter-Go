package controller

import (
	"net/http"

	"io/fs"

	"github.com/gin-gonic/gin"
)

// 主题列表处理函数
func ThemeListHandler(publicFS fs.FS) gin.HandlerFunc {
	return func(c *gin.Context) {
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
}
