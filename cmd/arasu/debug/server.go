package debug

import (
	"errors"
	"fmt"
	"github.com/arasuresearch/arasu/cmd/common"
	"log"
	"net"
	"net/http"
	"path"
	"strings"
	"time"
)

func (d *Debug) run() {
	if err := d.start(); err != nil {
		fmt.Println("Debug Server Starting  Error:" + err.Error())
		return
	}
	if d.DAS {
		http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path.Join(d.App.Root, "src/client/web/favicon.ico"))
		})
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if err := d.serve(w, r); err != nil {
				fmt.Println("Serving Error : " + err.Error())
				fmt.Fprintf(w, "%s", err.Error())
			}
		})
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				if _, ok := d.DGR[r.URL.Path]; !ok && d.serveIfAsset(w, r) {
					return
				}
				if err := d.serve(w, r); err != nil {
					fmt.Println("Serving Error : " + err.Error())
					fmt.Fprintf(w, "%s", err.Error())
				} else {
					d.DGR[r.URL.Path] = true
				}
				return
			}
			if err := d.serve(w, r); err != nil {
				fmt.Println("Serving Error : " + err.Error())
				fmt.Fprintf(w, "%s", err.Error())
			}
		})
	}
	fmt.Println("Server Started at : " + d.Port)
	//fmt.Println("you can browse now")
	log.Fatal(http.ListenAndServe(":"+d.Port, nil))

}

func (d *Debug) serveIfAsset(w http.ResponseWriter, r *http.Request) bool {
	defer r.Body.Close()
	r.Host = d.AssetUrl
	r.URL.Host = r.Host
	if len(r.URL.Scheme) == 0 {
		if r.Proto != "" {
			r.URL.Scheme = strings.ToLower(strings.Split(r.Proto, "/")[0])
		} else {
			r.URL.Scheme = "http"
		}
	}

	response, _ := d.Transport.RoundTrip(r)

	//fmt.Println(e0, response.Status, response.StatusCode)
	defer response.Body.Close()

	if response.StatusCode == 404 {
		return false
		//fmt.Fprintf(w, "Thani Hello, %q", html.EscapeString(r.URL.Path))
	} else {
		h, _ := w.(http.Hijacker)
		//fmt.Println(ok)
		conn, _, _ := h.Hijack()
		defer conn.Close()
		if err := response.Write(conn); err != nil {
			fmt.Fprintf(w, "%q", err.Error())
		}
	}
	return true
}

func (d *Debug) ignite() {
	l, _ := net.Listen("tcp", ":0")
	l.Close()
	_, d.DevPort, _ = net.SplitHostPort(l.Addr().String())
	d.Exe = path.Join(d.App.Root, "bin/debug") + " -p " + d.DevPort

	var err error
	d.Builder, err = d.App.NewBuilder()
	if err != nil {
		fmt.Println("new app building error", err)
		return
	}
	if !d.DAS {
		if len(d.AssetPort) == 0 {
			l, _ := net.Listen("tcp", ":0")
			l.Close()
			_, d.AssetPort, _ = net.SplitHostPort(l.Addr().String())
		}
		d.AssetUrl = d.AssetHost + ":" + d.AssetPort
		err = common.StartAssetServer(path.Join(d.App.Root, "src/client"), d.AssetMode, d.AssetPort)
		if err != nil {
			fmt.Println("starting asset error:" + err.Error())
			return
		}
		d.Transport = new(http.Transport)
		d.DGR = make(map[string]bool)
	}
	d.run()
}

func (d *Debug) start() error {
	//fmt.Println(d.Exe)
	d.Cmd = d.App.Cmd(d.Exe)
	if err := d.Cmd.Start(); err == nil {
		for {
			time.Sleep(100 * time.Millisecond)
			if _, e := net.Dial("tcp", ":"+d.DevPort); e == nil {
				break
			}
		}
		fmt.Println("debug server refreshed at localhost:" + d.DevPort)
	} else {
		return fmt.Errorf("App start Error : " + err.Error())
	}
	d.StartedAt = time.Now()
	return nil
}

func (d *Debug) stop() error {
	if err := d.Cmd.Process.Kill(); err == nil {
		//fmt.Println("sucess to kill: ", err)
		if err0 := d.Cmd.Process.Release(); err0 == nil {
			for {
				time.Sleep(200 * time.Millisecond)
				if _, e := net.Dial("tcp", ":"+d.DevPort); e != nil {
					break
				}
			}
			//fmt.Println("success to Releae: ", err0)
		} else {
			return errors.New("failed to Releae: " + err0.Error())
		}
	} else {
		return errors.New("failed to kill: " + err.Error())
	}
	return nil
}

func (d *Debug) ensure() error {
	if d.Builder.Watcher.NoChange() {
		if d.LastBuildError != nil {
			return d.LastBuildError
		}
		return nil
	}
	if err := d.Builder.ReBuild(); err != nil {
		d.LastBuildError = err
		return errors.New("Re Build" + err.Error())
	}
	d.LastBuildError = nil

	if err := d.stop(); err != nil {
		return errors.New("Stop Debug Server" + err.Error())
	}
	if err := d.start(); err != nil {
		return errors.New("Start Debug Server" + err.Error())
	}
	return nil

}

func (d *Debug) serve(w http.ResponseWriter, r *http.Request) error {
	if err := d.ensure(); err != nil {
		return err
	}

	defer r.Body.Close()
	transport := new(http.Transport)
	h, ok := w.(http.Hijacker)
	if !ok {
		return errors.New("Error on Transporting Connection")
	}

	r.Host = d.Host + ":" + d.DevPort
	//fmt.Println("Serving to dynamic server --> ", r.Host, r.URL.String())

	r.URL.Host = r.Host
	if len(r.URL.Scheme) == 0 {
		r.URL.Scheme = "http"
	}
	response, response_err := transport.RoundTrip(r)
	if response_err != nil {
		return errors.New("Error on Round Tripping" + response_err.Error())
	}
	connection, _, connection_err := h.Hijack()
	if connection_err != nil {
		return errors.New("Error on Hijacking  Connection" + connection_err.Error())
	}
	defer connection.Close()
	defer response.Body.Close()
	if err := response.Write(connection); err != nil {
		return errors.New("Writing Transport response on hijacked connection" + err.Error())
	}
	return nil
}
