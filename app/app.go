// --
// The BSD License (BSD)

// Copyright (c) 2015 Thaniyarasu Kannusamy <thaniyarasu@gmail.com> & Arasu Research Lab Pvt Ltd. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//    * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above copyright notice, this list of
//    conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//    * Neither Thaniyarasu Kannusamy <thaniyarasu@gmail.com>. nor ArasuResearch Inc may be used to endorse or promote products derived from this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND AUTHOR
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
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
