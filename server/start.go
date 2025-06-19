package server

import (
	"embed"
	"os"

	"github.com/spf13/cobra"
)

func StartWebServer(publicFS embed.FS) error {
	var rootCmd = &cobra.Command{
		Use:   "moeCounter",
		Short: "萌萌计数器！",
		Long:  "萌萌计数器是一个用于图片拼接的计数器服务，支持多种主题和配置选项。",
	}

	// 添加serve命令
	var port int
	var dbFile string
	var serveCmd = &cobra.Command{
		Use:   "start",
		Short: "启动Web服务",
		Run: func(cmd *cobra.Command, args []string) {
			// 启动并注册路由，传递嵌入的文件系统
			RegisterRoutes(port, dbFile, publicFS)
		},
	}
	serveCmd.Flags().IntVarP(&port, "port", "p", 8088, "服务监听端口")
	serveCmd.Flags().StringVarP(&dbFile, "db", "d", "data.db", "数据库文件路径")

	rootCmd.AddCommand(serveCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	return nil
}
