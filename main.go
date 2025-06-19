package main

import (
	"embed"
	"moeCounter/server"
)

//go:embed public
var publicFS embed.FS

func main() {
	server.StartWebServer(publicFS)
}
