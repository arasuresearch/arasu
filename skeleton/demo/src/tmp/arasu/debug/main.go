package main

import (
	"boot"
	"github.com/arasuresearch/arasu/cmd/arasu"
	_ "server/routes"
	_ "tmp/dispatchers"
)

func main() {
	arasu.StartDebugServer(boot.App)
}
