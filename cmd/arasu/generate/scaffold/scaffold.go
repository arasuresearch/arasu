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

package scaffold

import (
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/common"
	"github.com/arasuresearch/arasu/lib"
	"github.com/arasuresearch/arasu/lib/stringer"

	"path"
)

type Scaffold struct {
	common.SubCmd
	Dstore               string
	Name                 string
	FileName             string
	Cname                string
	Attrs                map[string]string
	ClientAttrs          interface{}
	ClientModelViewAttrs interface{}
	ParseArgs            []string

	ClientModelMetadata string

	GenArgs          []string
	MigrationVersion string
	SkeletonDir      string
	AppSrcDir        string
	Force            bool
	BigData          bool
}

func (c *Scaffold) Run() {
	if c.Help {
		fmt.Println(help_msg)
		return
	}
	if len(c.Args) < 2 || c.Help {
		fmt.Println(help_msg)
		return
	}

	// c.Name = lib.Camelize(c.Args[0])
	// c.Cname = lib.Constantize(c.Name)
	c.Name = stringer.Underscore(c.Args[0])
	c.Cname = stringer.Camelize(c.Name)

	c.Args = c.Args[1:]
	c.ParseArgs = c.ParseArgs[1:]
	// c.Attrs, c.GenArgs = lib.ParseKeyValueAndRemaningArguments(c.Args[1:])
	// for k, v := range c.Attrs {
	// 	c.Attrs[lib.Camelize(k)] = lib.Camelize(v)
	// }
	//fmt.Println(c.Args, c.Dstore)
	var err error
	if c.Dstore == "rdbms" {
		err = c.RdbmsRun()
	} else {
		err = c.BigdataRun()
	}
	fmt.Println("scaffold error", err)
}

func (c *Scaffold) Init(a *app.App, args []string) {
	c.App = a
	c.Args = args
	c.ParseArgs = lib.ArgsWithoutOptions(c.Args)
	//fmt.Println(args)
	c.Parse()
	c.MigrationVersion = lib.UniqueTimeVersion()
	c.SkeletonDir = path.Join(c.App.ArasuRoot, "skeleton/cmd/arasu/generate/scaffold")
	c.AppSrcDir = path.Join(c.App.Root, "src")
}

func (c *Scaffold) Parse() {
	c.Flag.BoolVar(&c.Help, "h", false, "a bool")
	c.Flag.BoolVar(&c.Help, "help", c.Help, "a bool")
	c.Flag.BoolVar(&c.Force, "f", false, "a bool")
	c.Flag.BoolVar(&c.Force, "force", c.Force, "a bool")
	c.Flag.StringVar(&c.Dstore, "ds", c.App.Conf["Dstore"], "Data store name")
	c.Flag.StringVar(&c.Dstore, "dstore", c.Dstore, "Data store name")
	c.Flag.Parse(c.App.FlagArgs)
}

var help_msg = `Usage:
  arasu generate scaffold NAME [field[:type][:index] field[:type][:index]] [options]

Options:
      [--skip-namespace]                        # Skip namespace (affects only isolated applications)
  -o, --orm=NAME                                # Orm to be invoked
                                                # Default: active_record
      [--force-plural]                          # Forces the use of a plural ModelName
      [--resource-route]                        # Indicates when to generate resource route
                                                # Default: true
  -y, [--stylesheets]                           # Generate Stylesheets
                                                # Default: true
  -se, [--stylesheet-engine=STYLESHEET_ENGINE]  # Engine for Stylesheets
                                                # Default: scss
      [--assets]                                # Indicates when to generate assets
                                                # Default: true
  -c, --scaffold-controller=NAME                # Scaffold controller to be invoked
                                                # Default: scaffold_controller

ActiveRecord options:
      [--migration]            # Indicates when to generate migration
                               # Default: true
      [--timestamps]           # Indicates when to generate timestamps
                               # Default: true
      [--parent=PARENT]        # The parent class for the generated model
      [--indexes]              # Add indexes for references and belongs_to columns
                               # Default: true
  -t, [--test-framework=NAME]  # Test framework to be invoked
                               # Default: test_unit

TestUnit options:
      [--fixture]                   # Indicates when to generate fixture
                                    # Default: true
  -r, [--fixture-replacement=NAME]  # Fixture replacement to be invoked

ScaffoldController options:
  -e, [--template-engine=NAME]  # Template engine to be invoked
                                # Default: erb
      [--helper]                # Indicates when to generate helper
                                # Default: true
      [--jbuilder]              # Indicates when to generate jbuilder
                                # Default: true

Asset options:
  -j, [--javascripts]                           # Generate JavaScripts
                                                # Default: true
  -je, [--javascript-engine=JAVASCRIPT_ENGINE]  # Engine for JavaScripts
                                                # Default: coffee

Runtime options:
  -f, [--force]    # Overwrite files that already exist
  -p, [--pretend]  # Run but do not make any changes
  -q, [--quiet]    # Suppress status output
  -s, [--skip]     # Skip files that already exist

Description:
    Scaffolds an entire resource, from model and migration to controller and
    views, along with a full test suite. The resource is ready to use as a
    starting point for your RESTful, resource-oriented application.

    Pass the name of the model (in singular form), either CamelCased or
    under_scored, as the first argument, and an optional list of attribute
    pairs.

    Attributes are field arguments specifying the model's attributes. You can
    optionally pass the type and an index to each field. For instance:
    "title body:text tracking_id:integer:uniq" will generate a title field of
    string type, a body with text type and a tracking_id as an integer with an
    unique index. "index" could also be given instead of "uniq" if one desires
    a non unique index.

    Timestamps are added by default, so you don't have to specify them by hand
    as 'created_at:datetime updated_at:datetime'.

    You don't have to think up every attribute up front, but it helps to
    sketch out a few so you can start working with the resource immediately.

    For example, 'scaffold post title body:text published:boolean' gives
    you a model with those three attributes, a controller that handles
    the create/show/update/destroy, forms to create and edit your posts, and
    an index that lists them all, as well as a resources :posts declaration
    in config/routes.rb.

    If you want to remove all the generated files, run
    'arasu destroy scaffold ModelName'.

Examples:
    arasu generate scaffold post
    arasu generate scaffold post title body:text published:boolean
    arasu generate scaffold purchase amount:decimal tracking_id:integer:uniq

`
