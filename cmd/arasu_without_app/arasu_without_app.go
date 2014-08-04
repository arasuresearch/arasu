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

package arasu_without_app

import (
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/arasu_without_app/new"
	"github.com/arasuresearch/arasu/cmd/common"
)

type ArasuWithoutApp struct {
	New new.New
	common.SubCmd
	Version bool
}

func (c *ArasuWithoutApp) Run() {
	if len(c.Args) == 0 {
		c.Help = true
	}
	if c.Help {
		fmt.Println(help_msg)
		return
	}
	if c.Version {
		fmt.Println(app.Runtime.Version)
		return
	}

	switch c.Args[0] {
	case "version":
		fmt.Println(app.Runtime.Version)
	case "help":
		fmt.Println(help_msg)
	default:
		fmt.Println(c.Args[0] + " is not a recognised command or this is not a arasu app")

	}
}
func (a *ArasuWithoutApp) Init(args []string) {
	a.Args = args
	a.Parse()
}
func (a *ArasuWithoutApp) Parse() {
	a.Flag.BoolVar(&a.Help, "h", false, "a bool")
	a.Flag.BoolVar(&a.Help, "help", a.Help, "a bool")
	a.Flag.BoolVar(&a.Version, "v", false, "a bool")
	a.Flag.BoolVar(&a.Version, "version", a.Version, "a bool")
	a.Flag.Parse(app.Runtime.FlagArgs)
}

var help_msg = `Usage:
  arasu new APP_PATH [options]

Options:
  -O, [--skip-active-record]     # Skip Active Record files
  -d, [--database=DATABASE]      # Preconfigure for selected database (options: mysql/oracle/postgresql/sqlite3/frontbase/ibm_db/sqlserver/jdbcmysql/jdbcsqlite3/jdbcpostgresql/jdbc)
                                 # Default: sqlite3

Runtime options:
  -f, [--force]    # Overwrite files that already exist
  -p, [--pretend]  # Run but do not make any changes
  -q, [--quiet]    # Suppress status output
  -s, [--skip]     # Skip files that already exist

Arasu options:
  -h, [--help]     # Show this help message and quit
  -v, [--version]  # Show Rails version number and quit

Description:
    The 'arasu new' command creates a new Arasu application with a default
    directory structure and configuration at the path you specify.

    You can specify extra command-line arguments to be used every time
    'arasu new' runs in the .arasurc configuration file in your home directory.

    Note that the arguments specified in the .arasurc file don't affect the
    defaults values shown above in this help message.

Example:
    arasu new ~/projects/go/demo

    This generates a skeletal Arasu installation in ~/project/go/demo.
    See the README in the newly created application to get going.
`
