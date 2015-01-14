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

package rdbms

import (
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/common"
	rd "github.com/arasuresearch/arasu/datastorage/rdbms"
	"path"
	"strings"
)

type Rdbms struct {
	common.SubCmd
	Mode   string
	DsPath string
}

func (c *Rdbms) Run() {
	if c.Help {
		fmt.Println(help_msg)
		return
	}

	schemaTemplatePath := path.Join(c.App.ArasuRoot, "skeleton/templates/schema.go.tmpl")
	modelRoot := path.Join(c.App.Root, "src/server/dstores", c.DsPath)
	adapterName, conf := c.App.DsNameConf(c.DsPath)
	args := strings.Join(c.Args, ":")

	migrater := rd.DstoreCommand{
		ModelRoot:          modelRoot,
		SchemaTemplatePath: schemaTemplatePath,
		AdapterName:        adapterName,
		Conf:               conf,
		Args:               args,
	}
	if err := migrater.Start(); err != nil {
		fmt.Println("Migration Error", err)
	}
}

func (c *Rdbms) Init(a *app.App, args []string) {
	c.App = a
	c.Args = args
	c.Parse()
}
func (c *Rdbms) Parse() {
	c.Flag.BoolVar(&c.Help, "h", false, "a bool")
	c.Flag.BoolVar(&c.Help, "help", c.Help, "a bool")
	c.Flag.StringVar(&c.Mode, "m", "dev", "a mode")
	c.Flag.StringVar(&c.Mode, "mode", c.Mode, "a mode")
	c.Flag.StringVar(&c.DsPath, "dp", "rdbms", "a bool")
	c.Flag.StringVar(&c.DsPath, "--ds-path", c.DsPath, "a bool")
	c.Flag.Parse(c.App.FlagArgs)
}

var help_msg = `Usage:
  arasu ds rd create [args] [options]

General args:
  all     # generate all mode database

General options:
  -m, [--mode]     # generate particular mode database
  -h, [--help]     # Print generator's options and usage
  -p, [--pretend]  # Run but do not make any changes
  -f, [--force]    # Overwrite files that already exist
  -s, [--skip]     # Skip files that already exist
  -q, [--quiet]    # Suppress status output

Example:
    arasu ds rd create

    This generates the relational database .
`
