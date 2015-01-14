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
