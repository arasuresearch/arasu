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

package migration

import (
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/common"
	"github.com/arasuresearch/arasu/lib"
)

type Migration struct {
	common.SubCmd
	Dstore  string
	Name    string
	Attrs   map[string]string
	Version string
}

func (c *Migration) Run() {
	if c.Help {
		fmt.Println(help_msg)
		return
	}
}

func (c *Migration) Init(a *app.App, args []string) {
	c.Attrs, c.Args = lib.ParseKeyValueAndRemaningArguments(args)
	c.Version = lib.UniqueTimeVersion()
	c.App = a
	c.Parse()
}

func (c *Migration) Parse() {
	c.Flag.BoolVar(&c.Help, "h", false, "a bool")
	c.Flag.BoolVar(&c.Help, "help", c.Help, "a bool")
	c.Flag.StringVar(&c.Dstore, "d", "rdbms", "Data store name")
	c.Flag.StringVar(&c.Dstore, "dstore", c.Dstore, "Data store name")
	c.Flag.Parse(c.App.FlagArgs)
}

var help_msg = `Usage:
  arasu generate migration name [args] [options]

General options:
  -d, [--dstore]     # generate particular mode database
  -h, [--help]       # Print generator's options and usage
  -p, [--pretend]  # Run but do not make any changes
  -f, [--force]    # Overwrite files that already exist
  -s, [--skip]     # Skip files that already exist
  -q, [--quiet]    # Suppress status output

Example:
    arasu generate migration add_name_to_user 

    This generates a migration called add_name_to_user .
`
