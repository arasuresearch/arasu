// --
// Copyright (c) 2014 Thaniyarasu Kannusamy <thaniyarasu@gmail.com> & Arasu Research , Inc.

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

package main

import (
	_ "bitbucket.org/pkg/inflect"
	_ "git.apache.org/thrift.git/lib/go/thrift"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/arasu_without_app"
	"github.com/arasuresearch/arasu/lib"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"

	"log"
	"reflect"
	"strings"
)

func main() {
	if app.Ok {
		app.Runtime.Build()
		if len(app.Runtime.Args) > 0 && app.Runtime.Args[0] == "serve" {
			mode := lib.ParseFlag([]string{"m", "mode"}, "debug", app.Runtime.Args)
			app.Runtime.Args[0] = mode
		}
		//fmt.Println(app.Runtime.CmdBin, app.Runtime.Args)
		cmdLine := app.Runtime.CmdBin + " " + strings.Join(app.Runtime.Args, " ")
		err := app.Runtime.Cmd(cmdLine).Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		args := app.Runtime.Args
		awa := struct {
			ArasuWithoutApp arasu_without_app.ArasuWithoutApp
		}{}
		cmd := reflect.ValueOf(&awa).Elem().Field(0)
		index := 0
		for _, e := range args {
			if c := cmd.FieldByName(strings.Title(e)); c.IsValid() {
				cmd = c
				index++
			} else {
				break
			}
		}
		args = args[index:]
		cmd = reflect.New(cmd.Type())
		cmd.MethodByName("Init").Call(app.Rvs{reflect.ValueOf(args)})
		cmd.MethodByName("Run").Call(app.Rvs{})
	}
}
