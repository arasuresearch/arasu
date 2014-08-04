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

package new

import (
	"errors"
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/common"
	"github.com/arasuresearch/arasu/lib"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type New struct {
	common.SubCmd
	Dstore   string
	Database string
	Name     string
}

func (n *New) Run() {
	if app.Runtime.Ok {
		msg, _ := lib.ParseAndExecuteTemplateText(help_msg_for_existing_app_dir, lib.H{"AppDir": app.Runtime.Root})
		fmt.Println(msg)
		return
	}

	if len(n.Args) == 0 {
		fmt.Println("application name can't be empty, Options may  be given after the application name. For details run: arasu new --help")
		return
	}
	if n.Help {
		fmt.Println(help_msg)
		return
	}
	dir := n.Args[0]

	if ok, _ := lib.IsExist(dir); ok {
		fmt.Println("directory already exists:", dir)
		return
	}
	n.Name = dir

	SkeletonDemoDir := path.Join(app.Runtime.ArasuRoot, "skeleton/demo")
	err := filepath.Walk(SkeletonDemoDir, func(src string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			dst := path.Join(dir, strings.TrimPrefix(src, SkeletonDemoDir))
			fmt.Println("create ", dst)

			if err := lib.CreateTemplatedFile(src, dst, n); err != nil {
				fmt.Println(err)
				return errors.New("On Creating New Project " + err.Error())
			}

			// if _, err := lib.CopyFile(src, dst); err != nil {
			// 	fmt.Println(err)
			// 	return errors.New("On Creating New Project " + err.Error())
			// }

		}
		return nil
	})
	if err != nil {
		fmt.Println("error occured ", err)
	}
	fmt.Println("New Arasu Project created at ", dir, " sucessfully")

	// pubGet := exec.Command("pub", "get")
	// pubGet.Stdin = os.Stdin
	// pubGet.Stdout = os.Stdout
	// pubGet.Stderr = os.Stderr
	// pubGet.Dir = path.Join(dir, "src/client")
	// if err := pubGet.Start(); err != nil {
	// 	fmt.Println("dependencies not installed")
	// }
	// pubGet.Wait()

	// dstoreCreate := exec.Command("arasu", "dstore", n.Dstore, "create")
	// dstoreCreate.Stdin = os.Stdin
	// dstoreCreate.Stdout = os.Stdout
	// dstoreCreate.Stderr = os.Stderr
	// dstoreCreate.Dir = path.Join(dir)
	// if err := dstoreCreate.Start(); err != nil {
	// 	fmt.Println("App initialized error", err)
	// }
	// dstoreCreate.Wait()
	// if err == nil {
	// 	fmt.Println("if no error occured in above statements then , New Arasu Project created at ", dir, " sucessfully")
	// } else {
	// 	fmt.Println(err)
	// }
}

func (n *New) Init(args []string) {
	n.Args = args
	n.Parse()
}

func (n *New) Parse() {
	n.Flag.StringVar(&n.Dstore, "ds", "rdbms", "a string")
	n.Flag.StringVar(&n.Dstore, "dstore", n.Dstore, "a string")

	n.Flag.StringVar(&n.Database, "d", "mysql", "a string")
	n.Flag.StringVar(&n.Database, "database", n.Database, "a string")

	n.Flag.BoolVar(&n.Help, "h", false, "a bool")
	n.Flag.BoolVar(&n.Help, "help", n.Help, "a bool")
	n.Flag.Parse(app.Runtime.FlagArgs)
}

// var help_msg = "ArasuHelp"
// var help_msg_for_app = "ArasuHelp"
var help_msg_for_existing_app_dir = `Already '{{.AppDir}}' Contains an Arasu Application. 
So ,you can't initialize a new Arasu application within the directory of another , 
please create arasu app in a non-Arasu directory or sub directory.
Type 'arasu' for help.
`

var help_msg = `Usage:
  rails new APP_PATH [options]

Options:
  -r, [--ruby=PATH]              # Path to the Ruby binary of your choice
                                 # Default: /home/dev/.rbenv/versions/2.1.0/bin/ruby
  -m, [--template=TEMPLATE]      # Path to some application template (can be a filesystem path or URL)
      [--skip-gemfile]           # Don't create a Gemfile
  -B, [--skip-bundle]            # Don't run bundle install
  -G, [--skip-git]               # Skip .gitignore file
      [--skip-keeps]             # Skip source control .keep files
  -O, [--skip-active-record]     # Skip Active Record files
  -S, [--skip-sprockets]         # Skip Sprockets files
  -d, [--database=DATABASE]      # Preconfigure for selected database (options: mysql/oracle/postgresql/sqlite3/frontbase/ibm_db/sqlserver/jdbcmysql/jdbcsqlite3/jdbcpostgresql/jdbc)
                                 # Default: sqlite3
  -j, [--javascript=JAVASCRIPT]  # Preconfigure for selected JavaScript library
                                 # Default: jquery
  -J, [--skip-javascript]        # Skip JavaScript files
      [--dev]                    # Setup the application with Gemfile pointing to your Rails checkout
      [--edge]                   # Setup the application with Gemfile pointing to Rails repository
  -T, [--skip-test-unit]         # Skip Test::Unit files
      [--rc=RC]                  # Path to file containing extra configuration options for rails command
      [--no-rc]                  # Skip loading of extra configuration options from .railsrc file

Runtime options:
  -f, [--force]    # Overwrite files that already exist
  -p, [--pretend]  # Run but do not make any changes
  -q, [--quiet]    # Suppress status output
  -s, [--skip]     # Skip files that already exist

Rails options:
  -h, [--help]     # Show this help message and quit
  -v, [--version]  # Show Rails version number and quit

Description:
    The 'rails new' command creates a new Rails application with a default
    directory structure and configuration at the path you specify.

    You can specify extra command-line arguments to be used every time
    'rails new' runs in the .railsrc configuration file in your home directory.

    Note that the arguments specified in the .railsrc file don't affect the
    defaults values shown above in this help message.

Example:
    rails new ~/Code/Ruby/weblog

    This generates a skeletal Rails installation in ~/Code/Ruby/weblog.
    See the README in the newly created application to get going.
`
