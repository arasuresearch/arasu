package app

import (
	"github.com/arasuresearch/arasu/router"
	"os"
	"path"
	"strings"
)

func init() {
	Runtime.Pwd, _ = os.Getwd()
	Runtime.Root = app_dir(Runtime.Pwd)
	Runtime.Args = os.Args[1:]

	for i, e := range Runtime.Args {
		if strings.HasPrefix(e, "-") {
			Runtime.FlagArgs = Runtime.Args[i:]
			break
		}
	}

	Runtime.ArasuRoot = arasu_dir()
	Runtime.GoRoot = go_dir()
	Runtime.Version = version()
	if len(Runtime.Root) == 0 {
		return
	}
	Ok = true
	Runtime.Ok = Ok
	Runtime.Name = path.Base(Runtime.Root)
	Runtime.CmdSrcRoot = path.Join(Runtime.Root, "src/tmp/arasu")
	Runtime.CmdBinRoot = path.Join(Runtime.Root, "bin")
	Runtime.CmdBin = path.Join(Runtime.Root, "bin/arasu")

	Runtime.GoCommands = go_commands()
	Runtime.Env = app_env()
	Runtime.Routes = &router.Router{&router.Route{}}
	Runtime.Registry = make(router.Registry)

}
func default_conf() map[string]string {
	return map[string]string{
		"Mode": "debug",
		"Host": "localhost",
		"Port": "4000",
		"DAS":  "false",
		//"AssetMode": "debug",
		//"AssetHost": "localhost",
	}
}
