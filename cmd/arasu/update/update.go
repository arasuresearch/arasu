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
