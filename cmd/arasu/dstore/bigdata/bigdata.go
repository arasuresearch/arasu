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

package bigdata

import (
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/common"
	bd "github.com/arasuresearch/arasu/datastorage/bigdata"
	"path"
	"strings"
)

type Bigdata struct {
	Mode string
	common.SubCmd
	DsPath string
}

func (c *Bigdata) Run() {
	if c.Help {
		fmt.Println(help_msg)
		return
	}
	schemaTemplatePath := path.Join(c.App.ArasuRoot, "skeleton/templates/schema.go.tmpl")
	modelRoot := path.Join(c.App.Root, "src/server/dstores", c.DsPath)
	adapterName, conf := c.App.DsNameConf(c.DsPath)
	args := strings.Join(c.Args, ":")

	migrater := bd.DstoreCommand{
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
func (c *Bigdata) Init(a *app.App, args []string) {
	c.App = a
	c.Args = args
	c.Parse()
}
func (c *Bigdata) Parse() {
	c.Flag.BoolVar(&c.Help, "h", false, "a bool")
	c.Flag.BoolVar(&c.Help, "help", c.Help, "a bool")
	c.Flag.StringVar(&c.Mode, "m", "dev", "a mode")
	c.Flag.StringVar(&c.Mode, "mode", c.Mode, "a mode")
	c.Flag.StringVar(&c.DsPath, "dp", "bigdata", "a bool")
	c.Flag.StringVar(&c.DsPath, "--ds-path", c.DsPath, "a bool")
	c.Flag.Parse(c.App.FlagArgs)
}

var help_msg = `Usage:
  arasu ds bd create [args] [options]

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
    arasu ds bd create

    This generates the relational database .
`
