// --
// Copyright (c) 2014 Thaniyarasu Kannusamy <thaniyarasu@gmail.com>.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
// ++

package app

import (
	"github.com/arasuresearch/arasu/router"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

var (
	Ok      bool
	Runtime = new(App)
)

type Rvs []reflect.Value
type App struct {
	Pwd        string
	Root       string
	Name       string
	CmdSrcRoot string
	CmdBinRoot string
	CmdBin     string

	ArasuRoot  string
	GoRoot     string
	Version    string
	GoCommands map[string]string
	Env        []string
	Ok         bool
	Mode       string
	AllConf    map[string]map[string]string
	Conf       map[string]string
	Args       []string
	FlagArgs   []string
	Routes     *router.Router
	Registry   router.Registry
}

func New() *App {
	a := new(App)
	src := reflect.ValueOf(Runtime).Elem()
	dst := reflect.ValueOf(a).Elem()
	for i := 0; i < src.NumField(); i++ {
		dst.Field(i).Set(src.Field(i))
	}
	return a
}

func (a *App) DsNameConf(src string) (string, string) {
	if v, ok := a.Conf[src]; ok {
		if a := strings.SplitN(v, ",", 2); len(a) == 2 {
			return a[0], a[1]
		}
	}
	return "", ""
}
func (a *App) Cmd(str string) *exec.Cmd {
	strs := strings.Split(str, " ")
	cmd := exec.Command(strs[0], strs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = a.Root
	cmd.Env = a.Env
	return cmd
}

// func (a *App) CmdWithConf(str string) *exec.Cmd {
// 	strs := strings.Split(str, " ")
// 	cmd := exec.Command(strs[0], strs[1:]...)
// 	cmd.Stdin = os.Stdin
// 	// cmd.Stdout = os.Stdout
// 	// cmd.Stderr = os.Stderr
// 	cmd.Dir = a.Root
// 	cmd.Env = a.Env
// 	return cmd
// }

// func NameConf(src string) (string, string, error) {
// 	if v, ok := Conf[src]; ok {
// 		if a := strings.SplitN(v, ",", 2); len(a) == 2 {
// 			return a[0], a[1], nil
// 		}
// 	}
// 	return "", "", fmt.Errorf("can't parse %s configuration ", src)
// }
