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

package debug

import (
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/common"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

type Debug struct {
	common.SubCmd

	Mode string
	Host string
	Port string

	DAS      bool
	AssetUrl string
	DevPort  string
	Exe      string

	//	AssetServer *app.AssetServer
	Builder   *app.Builder
	Transport *http.Transport
	StartedAt time.Time
	Cmd       *exec.Cmd
	DGR       map[string]bool //Dynamic Get Requests

	AssetMode string
	AssetHost string
	AssetPort string

	LastBuildError error
}

func (d *Debug) Run() {
	if d.Help {
		fmt.Println(help_msg)
		return
	}
	if _, e := net.Dial("tcp", ":"+d.Port); e == nil {
		fmt.Println(d.Port + " port is not free")
		return
	}
	d.ignite()
}

func (d *Debug) Init(ap *app.App, args []string) {
	d.App = ap
	d.Args = args
	d.Parse()
}

func (d *Debug) Parse() {
	f := d.Flag
	c := d.App.Conf

	mode, _ := c["Mode"]
	host, _ := c["Host"]
	port, _ := c["Port"]
	das, _ := strconv.ParseBool(c["DAS"])
	amode, _ := c["AssetMode"]
	ahost, _ := c["AssetHost"]
	aport, _ := c["AssetPort"]

	f.StringVar(&mode, "m", mode, "a string")
	f.StringVar(&mode, "mode", mode, "a string")
	f.StringVar(&port, "p", port, "a string")
	f.StringVar(&port, "port", port, "a string")
	f.StringVar(&host, "hos", host, "a string")
	f.StringVar(&host, "host", host, "a string")

	f.BoolVar(&das, "das", das, "a bool")
	f.BoolVar(&das, "disable-asset-server", das, "a bool")

	f.StringVar(&amode, "am", amode, "a string")
	f.StringVar(&amode, "asset-mode", amode, "a string")
	f.StringVar(&ahost, "ah", ahost, "a string")
	f.StringVar(&ahost, "asset-host", ahost, "a string")
	f.StringVar(&aport, "ap", aport, "a string")
	f.StringVar(&aport, "asset-port", aport, "a string")

	f.BoolVar(&d.Help, "h", false, "a bool")
	f.BoolVar(&d.Help, "help", d.Help, "a bool")
	f.Parse(d.App.FlagArgs)
	if len(amode) == 0 {
		amode = mode
	}
	d.Mode = mode
	d.Host = host
	d.Port = port
	d.DAS = das
	d.AssetMode = amode
	d.AssetHost = ahost
	d.AssetPort = aport
}

var help_msg = `Usage: arasu server [mongrel, thin, etc] [options]
    -p, --port=port                  Runs Arasu on the specified port.
                                     Default: 4000
    -b, --binding=ip                 Binds Arasu to the specified ip.
                                     Default: 0.0.0.0
    -c, --config=file                Use custom rackup configuration file
    -d, --daemon                     Make server run as a Daemon.
    -u, --debugger                   Enable the debugger
    -m, --mode=name                  Specifies the mode to run this server under (tes/dev/pro).
                                     Default: dev
    -P, --pid=pid                    Specifies the PID file.
                                     Default: tmp/pids/server.pid

    -h, --help                       Show this help message.
`
