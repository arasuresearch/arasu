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
