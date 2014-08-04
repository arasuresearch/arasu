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

package update

import (
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/arasu/update/schema"
	"github.com/arasuresearch/arasu/cmd/common"
	"github.com/arasuresearch/arasu/lib"
)

type Update struct {
	Schema schema.Schema
	common.SubCmd
}

func (c *Update) Run() {
	if len(c.Args) == 0 {
		fmt.Println(help_msg)
		return
	}
	if c.Help {
		fmt.Println(help_msg)
		return
	}
	msg, _ := lib.ParseAndExecuteTemplateText(help_msg_for_sub_command, lib.HSS{"Name": c.Args[0]})
	fmt.Println(msg)
}

func (c *Update) Init(a *app.App, args []string) {
	c.App = a
	c.Args = args
	c.Parse()
}
func (c *Update) Parse() {
	c.Flag.BoolVar(&c.Help, "h", false, "a bool")
	c.Flag.BoolVar(&c.Help, "help", c.Help, "a bool")
	c.Flag.Parse(c.App.FlagArgs)
}

var help_msg_for_sub_command = `Could not find generator {{.Name}}.
Type 'arasu generate -h'  for help.
`
var help_msg = `Usage: arasu generate GENERATOR [args] [options]

General options:
  -h, [--help]     # Print generator's options and usage
  -p, [--pretend]  # Run but do not make any changes
  -f, [--force]    # Overwrite files that already exist
  -s, [--skip]     # Skip files that already exist
  -q, [--quiet]    # Suppress status output

Please choose a generator below.
  server:
    assets
    controller
    generator
    helper
    integration_test
    jbuilder
    mailer
    migration
    model
    resource
    scaffold
    scaffold_controller
    task
  client:
    assets
    controller
    generator
    helper
    integration_test
    jbuilder
    mailer
    migration
    model
    resource
    scaffold
    scaffold_controller
    task
`
