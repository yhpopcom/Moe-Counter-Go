package cmd

import (
	"moeCounter/cmd/flags"
	"moeCounter/public"
	"moeCounter/server"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "start",
	Short: "启动Web服务",
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化路由
		router := server.InitRouter(flags.Port, flags.DbFile, public.Public, flags.Debug)
		// 启动服务器
		server.RunServer(router, flags.Port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
