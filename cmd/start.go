package cmd

import (
	"moeCounter/public"
	"moeCounter/server"

	"github.com/spf13/cobra"
)

// 添加serve命令
var port int
var dbFile string
var debug bool
var serveCmd = &cobra.Command{
	Use:   "start",
	Short: "启动Web服务",
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化路由
		router := server.InitRouter(port, dbFile, public.Public, debug)
		// 启动服务器
		server.RunServer(router, port)
	},
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 8088, "服务监听端口")
	serveCmd.Flags().StringVarP(&dbFile, "db", "d", "data.db", "数据库文件路径")
	serveCmd.Flags().BoolVar(&debug, "debug", false, "是否开启调试模式")
	rootCmd.AddCommand(serveCmd)
}
