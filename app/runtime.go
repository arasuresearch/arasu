package app

import (
	"github.com/arasuresearch/arasu/lib"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func SetConf(conf map[string]map[string]string) {
	def := conf[""]
	for k, v := range default_conf() {
		if _, ok := def[k]; !ok {
			def[k] = v
		}
	}

	Runtime.Mode = lib.ParseFlag([]string{"m", "mode"}, "debug", Runtime.Args)

	Runtime.AllConf = conf
	Runtime.Conf = conf[""]
	for k, v := range conf[Runtime.Mode] {
		Runtime.Conf[k] = v
	}
}

func IsThisGoCommand(cmd string) bool {
	_, ok := Runtime.GoCommands[cmd]
	return ok
}
func arasu_dir() string {
	_, this_file, _, _ := runtime.Caller(0)
	return path.Dir(path.Dir(this_file))
}

//getting arasu present app dir from command line
func app_dir(pwd string) string {
	dir := pwd
	for dir != path.Dir(dir) {
		if _, err := os.Stat(path.Join(dir, ".arasu_project_rc")); err == nil {
			return dir
		}
		dir = path.Dir(dir)
	}
	return ""
}

func go_dir() string {
	out, _ := exec.Command("go", "env").Output()
	for _, e := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(e, "GOROOT=") {
			e = strings.TrimPrefix(e, "GOROOT=")
			e = strings.Trim(e, `"`)
			return e
		}
	}
	return ""
}

func app_env() []string {
	env := os.Environ()
	for i, e := range env {
		if strings.HasPrefix(e, "GOPATH=") {
			gopath := strings.Split(env[i], "=")
			gopath[1] = Runtime.Root + string(os.PathListSeparator) + gopath[1]
			env[i] = strings.Join(gopath, "=")
			break
		}
	}
	for i, e := range env {
		if strings.HasPrefix(e, "GOBIN=") {
			env[i] = ""
			break
		}
	}
	return env
}

func go_commands() map[string]string {
	commands := make(map[string]string)
	if out, err := exec.Command("go", "help").Output(); err == nil {
		temp := strings.Split(string(out), "The commands are:")[1]
		temp = strings.Split(temp, "Use \"go help [command]\"")[0]
		temp = strings.Trim(temp, "\n")
		for _, line := range strings.Split(temp, "\n") {
			l0 := strings.Split(strings.TrimSpace(line), " ")
			commands[l0[0]] = strings.TrimSpace(strings.Join(l0[1:], " "))
		}
	} else {
		log.Fatal(err)
	}
	return commands
}
func version() string {
	if data, err := ioutil.ReadFile(path.Join(Runtime.ArasuRoot, "VERSION")); err == nil {
		return string(data)
	} else {
		log.Fatal(err)
	}
	return ""
}
