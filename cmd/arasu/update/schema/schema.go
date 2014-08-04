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

package schema

import (
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/common"
	"github.com/arasuresearch/arasu/datastorage/bigdata"
	"github.com/arasuresearch/arasu/datastorage/rdbms"
	"path"
	"strings"
)

type Schema struct {
	common.SubCmd
	Dstore string
	Name   string
	Attrs  map[string]string
}

func (c *Schema) Run() {
	// if len(c.Args) == 0 {
	// 	fmt.Println(help_msg)
	// 	return
	// }
	if c.Help {
		fmt.Println(help_msg)
		return
	}
	if err := c.SchemaToStruct(); err != nil {
		fmt.Println(err)
	}

}

func (c *Schema) Init(a *app.App, args []string) {
	c.App = a
	c.Args = args
	c.Parse()
}

func (c *Schema) Parse() {
	c.Flag.BoolVar(&c.Help, "h", false, "a bool")
	c.Flag.BoolVar(&c.Help, "help", c.Help, "a bool")
	c.Flag.StringVar(&c.Dstore, "ds", c.App.Conf["Dstore"], "a bool")
	c.Flag.StringVar(&c.Dstore, "dstore", c.Dstore, "a bool")
	c.Flag.Parse(c.App.FlagArgs)
}

func (c *Schema) SchemaToStruct() error {
	templatePath := path.Join(c.App.ArasuRoot, "skeleton/templates/schema.go.tmpl")
	dsPath := c.Dstore
	// TODO remove slash later
	name, conf := c.App.DsNameConf(dsPath)
	var err error

	switch {
	case strings.HasPrefix(dsPath, "bigdata"):
		err = bigdata.SchemaToStruct(name, conf, path.Join(c.App.Root, "src/server/dstores", dsPath), templatePath)
	case strings.HasPrefix(dsPath, "rdbms"):
		fmt.Println(name, conf, dsPath)
		err = rdbms.SchemaToStruct(name, conf, path.Join(c.App.Root, "src/server/dstores", dsPath), templatePath)
	}
	if err != nil {
		fmt.Println("Error while update schema to struct ", err)

	}

	// err := filepath.Walk(dsRoot, func(src string, info os.FileInfo, err error) error {
	// 	if !info.IsDir() || src == dsRoot {
	// 		return nil
	// 	}
	// 	name := strings.TrimPrefix(src, dsRoot)
	// 	dsType, conf := c.App.DsNameConf(name)
	// 	switch {
	// 	case strings.HasPrefix(name, "/bigdata"):
	// 		bigdata.SchemaToStruct(dsType, conf, path.Join(dsRoot, name), templatePath)
	// 	case strings.HasPrefix(name, "/rdbms"):
	// 		rdbms.SchemaToStruct(dsType, conf, path.Join(dsRoot, name), templatePath)
	// 	default:
	// 		fmt.Println("can't identify ds type " + name)
	// 	}

	// 	return nil
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

var help_msg = `Usage:
  arasu schema update --dstore [rdbms|bigdata]
`
var help_msg_for_sub_command = `Could not find generator {{.Name}}.
Type 'arasu generate -h'  for help.
`
