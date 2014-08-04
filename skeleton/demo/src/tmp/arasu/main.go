package main

import (
	"boot"
	"github.com/arasuresearch/arasu/cmd/arasu"
	_ "server/routes"
	_ "tmp/dispatchers"
)

func main() {
	arasu.Start(boot.App)
}

// _ "ds/bigdata/migrate"
// _ "ds/rdbms/migrate"
