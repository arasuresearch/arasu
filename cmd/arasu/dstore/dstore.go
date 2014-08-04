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

package dstore

import (
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/common"
	"github.com/arasuresearch/arasu/datastorage/bigdata"
	"github.com/arasuresearch/arasu/datastorage/rdbms"
	"github.com/arasuresearch/arasu/lib"
	"path"
	"strings"
)

type Dstore struct {
	// Bigdata bigdata.Bigdata
	// Rdbms   rdbms.Rdbms
	common.SubCmd
	DsPath string
	Dstore string
}

func (c *Dstore) Run() {
	if c.Help {
		fmt.Println(help_msg)
		return
	}
	schemaTemplatePath := path.Join(c.App.ArasuRoot, "skeleton/templates/schema.go.tmpl")
	modelRoot := path.Join(c.App.Root, "src/server/dstores", c.DsPath)
	adapterName, conf := c.App.DsNameConf(c.DsPath)
	args := strings.Join(lib.ArgsWithoutOptions(c.Args), ":")
	//fmt.Println(c.DsPath)
	if strings.HasPrefix(c.DsPath, "rdbms") {
		migrater := rdbms.DstoreCommand{
			ModelRoot:          modelRoot,
			SchemaTemplatePath: schemaTemplatePath,
			AdapterName:        adapterName,
			Conf:               conf,
			Args:               args,
		}
		if err := migrater.Start(); err != nil {
			fmt.Println("Migration Error", err)
		}
	} else {
		migrater := bigdata.DstoreCommand{
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

}

func (c *Dstore) Init(a *app.App, args []string) {
	c.App = a
	c.Args = args
	c.Parse()
}
func (c *Dstore) Parse() {
	c.Flag.BoolVar(&c.Help, "h", false, "a bool")
	c.Flag.BoolVar(&c.Help, "help", c.Help, "a bool")
	c.Flag.StringVar(&c.Dstore, "d", c.App.Conf["Dstore"], "Data store name")
	c.Flag.StringVar(&c.Dstore, "ds", c.Dstore, "Data store name")
	c.Flag.StringVar(&c.Dstore, "dstore", c.Dstore, "Data store name")
	c.Flag.StringVar(&c.DsPath, "dp", "", "a bool")
	c.Flag.StringVar(&c.DsPath, "--ds-path", c.DsPath, "a bool")
	c.Flag.Parse(c.App.FlagArgs)
	if len(c.DsPath) == 0 {
		c.DsPath = c.Dstore
	}
}

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
