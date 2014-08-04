package boot

import (
	_ "boot/initializers"
	"config"
	_ "config/debug"
	_ "config/release"
	_ "config/testing"
	"github.com/arasuresearch/arasu/app"
)

var App *app.App

func init() {
	app.SetConf(config.Config)
	App = app.New()
}
