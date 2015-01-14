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

package main

//TODO : remove void import packages later
import (
	//_ "git.apache.org/thrift.git/lib/go/thrift"
	_ "github.com/apache/thrift/lib/go/thrift"
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
