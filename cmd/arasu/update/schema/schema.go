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
