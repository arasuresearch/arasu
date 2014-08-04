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

package release

import (
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/cmd/common"
	"github.com/arasuresearch/arasu/handler"
	"log"
	"net"
	"net/http"
	"path"
	"strconv"
)

type Release struct {
	common.SubCmd

	Mode string
	Host string
	Port string
	DAS  bool

	AssetMode string
	AssetHost string
	AssetPort string

	AssetUrl string
}
type ReleaseHandler struct {
	Transport *http.Transport
	DGR       map[string]bool //Dynamic Get Requests
	AssetUrl  string
}

func (r *Release) NewReleaseHandler() *handler.ReleaseHandler {
	rh := &handler.ReleaseHandler{
		App:      r.App,
		Registry: r.App.Registry,
		Routes:   r.App.Routes,
		DAS:      r.DAS,
	}
	if !r.DAS {
		if len(r.AssetPort) == 0 {
			l, _ := net.Listen("tcp", ":0")
			l.Close()
			_, r.AssetPort, _ = net.SplitHostPort(l.Addr().String())
		}
		rh.AssetUrl = r.AssetHost + ":" + r.AssetPort
		err := common.StartAssetServer(path.Join(r.App.Root, "src/client"), r.AssetMode, r.AssetPort)
		if err != nil {
			fmt.Println("starting asset error:" + err.Error())
			return rh
		}
		rh.DGR = make(map[string]bool)
		rh.Transport = new(http.Transport)
	}
	return rh
}

func (r *Release) Run() {
	if r.Help {
		fmt.Println(help_msg)
		return
	}
	if _, e := net.Dial("tcp", ":"+r.Port); e == nil {
		fmt.Println(r.Port + " port is not free")
		return
	}
	if err := r.App.Build(); err != nil {
		fmt.Println(err)
		return
	}
	hand := r.NewReleaseHandler()
	log.Fatal(http.ListenAndServe(":"+r.Port, hand))
}

func (r *Release) Init(ap *app.App, args []string) {
	r.App = ap
	r.Args = args
	r.Parse()
}

func (r *Release) Parse() {
	f := r.Flag
	c := r.App.Conf

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

	f.BoolVar(&r.Help, "h", false, "a bool")
	f.BoolVar(&r.Help, "help", r.Help, "a bool")
	f.Parse(r.App.FlagArgs)
	if len(amode) == 0 {
		amode = mode
	}
	if len(ahost) == 0 {
		ahost = host
	}
	r.Mode = mode
	r.Host = host
	r.Port = port
	r.DAS = das
	r.AssetMode = amode
	r.AssetHost = ahost
	r.AssetPort = aport
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
