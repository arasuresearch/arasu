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

package about

import (
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/common"
)

type About struct {
	common.SubCmd
}

func (a *About) Run() {
	if a.Help {
		fmt.Println(help_msg)
		return
	}
}
func (a *About) Init(ap *app.App, args []string) {
	a.App = ap
	a.Args = args
	a.Parse()
}
func (a *About) Parse() {
	a.Flag.BoolVar(&a.Help, "h", false, "a bool")
	a.Flag.BoolVar(&a.Help, "help", a.Help, "a bool")
	a.Flag.Parse(a.App.FlagArgs)
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
