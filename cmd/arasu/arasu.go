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

package arasu

import (
	"flag"
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/arasu/about"
	"github.com/arasuresearch/arasu/cmd/arasu/debug"
	"github.com/arasuresearch/arasu/cmd/arasu/dstore"
	"github.com/arasuresearch/arasu/cmd/arasu/generate"
	"github.com/arasuresearch/arasu/cmd/arasu/release"
	"github.com/arasuresearch/arasu/cmd/arasu/update"
	"github.com/arasuresearch/arasu/cmd/common"
	"github.com/arasuresearch/arasu/handler"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type Arasu struct {
	// Server   server.Server
	About    about.About
	Generate generate.Generate
	Dstore   dstore.Dstore
	Update   update.Update
	Debug    debug.Debug
	Release  release.Release
	common.SubCmd
}

func Start(a *app.App) {
	args := app.Runtime.Args
	cmd := reflect.ValueOf(&struct{ Arasu Arasu }{}).Elem().Field(0)
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
	rargs := app.Rvs{
		reflect.ValueOf(a),
		reflect.ValueOf(args),
	}
	cmd.MethodByName("Init").Call(rargs)
	cmd.MethodByName("Run").Call(app.Rvs{})
}

func StartDebugServer(a *app.App) {
	var host, port string
	f := flag.FlagSet{}
	f.StringVar(&port, "p", "4001", "a string")
	f.StringVar(&host, "h", "localhost", "a string")
	f.Parse(os.Args[1:])

	dh := &handler.DebugHandler{
		App:      a,
		Registry: a.Registry,
		Routes:   a.Routes,
	}
	log.Fatal(http.ListenAndServe(":"+port, dh))
}

func (a *Arasu) Run() {
	if a.Help || len(a.Args) == 0 {
		fmt.Println(help_msg)
		return
	}
	if app.IsThisGoCommand(a.Args[0]) {
		err := a.App.Cmd("go " + strings.Join(a.Args, " ")).Run()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Error: Command '" + a.Args[0] + "' not recognized")
		fmt.Println(help_msg)
	}
}
func (a *Arasu) Init(ap *app.App, args []string) {
	a.App = ap
	a.Args = args
	a.Parse()
}
func (a *Arasu) Parse() {
	a.Flag.BoolVar(&a.Help, "h", false, "a bool")
	a.Flag.BoolVar(&a.Help, "help", a.Help, "a bool")
	a.Flag.Parse(a.Args)
}

var help_msg = `Usage: arasu COMMAND [ARGS]

The most common arasu commands are:
 generate    Generate new code (short-cut alias: "g")
 console     Start the Arasu console (short-cut alias: "c")
 server      Start the Arasu server (short-cut alias: "s")
 dbconsole   Start a console for the database specified in config/database.yml
             (short-cut alias: "db")
 new         Create a new Arasu application. "arasu new myapp" creates a
             new application called myapp in "./myapp"

In addition to those, there are:
 application  Generate the Arasu application code
 destroy      Undo code generated with "generate" (short-cut alias: "d")
 plugin new   Generates skeleton for developing a Arasu plugin
 runner       Run a piece of code in the application environment (short-cut alias: "r")

All commands can be run with -h (or --help) for more information.
`
