package handler

import (
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/router"
	"net/http"
	"path"
	"strings"
)

type ReleaseHandler struct {
	App       *app.App
	Transport *http.Transport
	DGR       map[string]bool //Dynamic Get Requests
	AssetUrl  string
	DAS       bool //Disable Asset Server
	Routes    *router.Router
	Registry  router.Registry
}

func (h *ReleaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.DAS {
		if r.Method == "GET" && (r.URL.Path == "/" || r.URL.Path == "/favicon.ico") {
			fileName := "src/client/web/favicon.ico"
			if r.URL.Path == "/" {
				fileName = "src/client/web/index.html"
			}
			http.ServeFile(w, r, path.Join(h.App.Root, fileName))
			return
		}
		if r.Method == "GET" && r.URL.Path == "/favicon.ico" {
			http.ServeFile(w, r, path.Join(h.App.Root, "src/client/web/favicon.ico"))
			return
		}
		if err := handle(w, r, h.Registry, h.Routes, h.App); err != nil {
			fmt.Fprintf(w, "%s", err.Error())
		}
	} else {
		if r.Method == "GET" {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "index.html", 302)
				return
			}
			if _, ok := h.DGR[r.URL.Path]; !ok && h.serveIfAsset(w, r) {
				return
			}
			if err := handle(w, r, h.Registry, h.Routes, h.App); err != nil {
				fmt.Fprintf(w, "%s", err.Error())
			} else {
				h.DGR[r.URL.Path] = true
			}
		} else {
			if err := handle(w, r, h.Registry, h.Routes, h.App); err != nil {
				fmt.Fprintf(w, "%s", err.Error())
			}
		}

	}
}

func (h *ReleaseHandler) serveIfAsset(w http.ResponseWriter, r *http.Request) bool {
	defer r.Body.Close()
	r.Host = h.AssetUrl
	r.URL.Host = r.Host
	if len(r.URL.Scheme) == 0 {
		if r.Proto != "" {
			r.URL.Scheme = strings.ToLower(strings.Split(r.Proto, "/")[0])
		} else {
			r.URL.Scheme = "http"
		}
	}
	response, _ := h.Transport.RoundTrip(r)
	defer response.Body.Close()
	if response.StatusCode == 404 {
		return false
	} else {
		h, _ := w.(http.Hijacker)
		conn, _, _ := h.Hijack()
		defer conn.Close()
		if err := response.Write(conn); err != nil {
			fmt.Fprintf(w, "%q", err.Error())
		}
	}
	return true
}
