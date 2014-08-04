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

package router

import (
	"fmt"
	"github.com/arasuresearch/arasu/lib"
	"net/http"
	"strings"
)

func (r *Router) Set(path string, stringMatches [][]string) {
	matches := make([]Match, len(stringMatches))
	for i, e := range stringMatches {
		matches[i] = Match{e[0], e[1], e[2]}
	}
	paths := StringPathSplit(path)
	route := r.All
	for _, e := range paths {
		if route.Paths == nil {
			route.Paths = make(Path)
		}
		if _, ok := route.Paths[e]; !ok {
			route.Paths[e] = &Route{}
		}
		route = route.Paths[e]
	}
	route.Matches = append(route.Matches, matches...)
}

func (r *Router) SetMatches(path string, matches []Match) {
	paths := StringPathSplit(path)
	route := r.All
	for _, e := range paths {
		if route.Paths == nil {
			route.Paths = make(Path)
		}
		if _, ok := route.Paths[e]; !ok {
			route.Paths[e] = &Route{}
		}
		route = route.Paths[e]
	}
	route.Matches = append(route.Matches, matches...)
}

func (rr *Router) Find(w *http.ResponseWriter, r *http.Request, registry Registry) (cntr Handler, cname string, fname string, format string, params map[string]string, err error) {
	// fmt.Println("================")
	// fmt.Println(r.RequestURI, r.URL, r.URL.String(), r.URL.Path)

	// fmt.Println("================")

	paths := StringPathSplit(r.URL.Path)
	format = ResolveFormat(r)
	route := rr.All
	if len(paths) == 0 {
		cntr, cname, fname, params, err = FindMatches(route.Matches, r.Method, paths, 0, format, registry)
		return
	}
	for i, e := range paths {
		lroute, ok := route.Paths[e]
		if !ok {
			lroute, ok = route.Paths["*"]
		}
		if !ok {
			if err == nil {
				args := []interface{}{r.URL.Path, strings.Join(paths[:i], "/"), strings.Join(paths[i:], "/")}
				err = fmt.Errorf("Error: No Path Found for %s ( %s - %s )", args...)
			}
			return
		}
		cntr, cname, fname, params, err = FindMatches(lroute.Matches, r.Method, paths, i+1, format, registry)
		if err == nil {
			return
		} else {
			route = lroute
		}
	}
	return
}
func FindMatches(matches []Match, method string, paths []string, i int, format string, registry Registry) (cntr Handler, cname string, fname string, params map[string]string, err error) {
	PostPath := strings.Join(paths[i:], "/")
	var ok, found bool
	for _, e := range matches {
		path := e.Path
		if strings.Contains(e.Path, ".") {
			lv := strings.Split(e.Path, ".")
			path = lv[0]
			if !strings.Contains(lv[1], format) {
				continue
			}
		}
		if (path == "*" || path == PostPath) && (e.Method == "*" || e.Method == method) {
			found = true
		}
		if e.Method == method && strings.Contains(path, ":") && strings.Count(path, "/") == len(paths)-(i+1) {
			if params, ok = RegExpMatch(PostPath, path); ok {
				found = true
			}
		}
		if found {
			match := AbsMatch(e, paths, i)
			cf := strings.Split(match.Action, ".")
			cname = cf[0]
			fname = cf[1]
			cntr, ok = registry["/"+cname]
			if ok {
				if !lib.StringArrayContains(cntr.Methods, fname) {
					err = fmt.Errorf("Action %s Not Found on Controller(%s) ", fname, cname)
				}
			} else {
				err = fmt.Errorf("Controller %s Not Found", cname)
			}
			return
		}
	}
	err = fmt.Errorf("Error:  No Match Found for %s", strings.Join(paths, "/"))
	return
}

func AbsMatch(m Match, paths []string, i int) Match {
	cname := strings.Join(paths[:i], "/")
	fname := strings.Title(firstValidPathPiece(paths[i:]))
	if fname == "" {
		fname = "Index"
	}

	action := m.Action
	if !strings.Contains(action, ".") {
		action = "." + action
	}

	actions := strings.Split(action, ".")
	if actions[0] == "" || actions[0] == "*" {
		actions[0] = cname
	}
	if actions[1] == "" || actions[1] == "*" {
		actions[1] = fname
	}

	m.Action = strings.Join(actions, ".")
	return m
}
func firstValidPathPiece(post_paths []string) string {
	for _, e := range post_paths {
		if !strings.Contains(e, ":") {
			return e
		}
	}
	return ""
}

func RegExpMatch(path string, match string) (map[string]string, bool) {
	matches := strings.Split(match, "/")
	result := make(map[string]string)
	for i, e := range strings.Split(path, "/") {
		if strings.HasPrefix(matches[i], ":") {
			result[matches[i][1:]] = e
		} else if matches[i] != e {
			return nil, false
		}
	}
	return result, true
}
func StringPathSplit(s string) (r []string) {
	for _, e := range strings.Split(strings.Split(s, ".")[0], "/") {
		if len(e) > 0 {
			r = append(r, e)
		}
	}
	return
}
func splitPath(s string) []string {
	a := strings.Split(s, "/")

	if len(a) > 0 && strings.Contains(a[len(a)-1], ".") {
		a[len(a)-1] = strings.Split(a[len(a)-1], ".")[0]
	}
	if len(a) > 0 && a[len(a)-1] == "" {
		a = a[:len(a)-1]
	}
	if len(a) > 0 && a[0] == "" {
		a = a[1:]
	}
	switch len(a) {
	case 0:
		a = []string{"", ""}
	case 1:
		a = append(a, "")
	}
	return a
}

func ResolveFormat(r *http.Request) string {
	if strings.Contains(r.URL.Path, ".") {
		t := strings.Split(r.URL.Path, ".")
		return t[len(t)-1]
	}

	text := r.Header.Get("accept")
	switch {
	case strings.Contains(text, "application/json"),
		strings.Contains(text, "text/javascript"),
		strings.Contains(text, "*/*"):
		return "json"
	case strings.Contains(text, "application/xml"),
		strings.Contains(text, "text/xml"):
		return "html"
	case strings.Contains(text, "text/plain"):
		return "text"
	}
	return "html"
}
