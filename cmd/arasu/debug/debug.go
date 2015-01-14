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
