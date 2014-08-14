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

package controller

import (
	"github.com/arasuresearch/arasu/app"
	"net/http"
)

type Controller struct {
	W         *http.ResponseWriter
	R         *http.Request
	App       *app.App
	Name      string
	Action    string
	Methods   []string
	FuncsArgs map[string][]string
	Params    map[string]interface{}
}

type BaseController struct {
	Controller
}

//////////////////////////////////
// package controller

// import (
// 	"arasu/app"
// 	"encoding/json"
// 	//"encoding"

// 	"fmt"
// 	"net/http"
// 	"path"
// 	"reflect"
// 	"strings"
// )

// type Controller struct {
// 	W   *http.ResponseWriter
// 	R   *http.Request
// 	App *app.App

// 	Name   string
// 	Action string
// 	Format string
// 	View   string
// 	Layout string

// 	Params  map[string]string
// 	Data    map[string]interface{}
// 	Methods []string
// 	Cntr    reflect.Value

// 	SkipLayout bool
// 	Rendered   bool
// 	Redirected bool
// 	Denied     bool
// }

// type BaseController struct {
// 	Controller
// }

// func (c *BaseController) Serve(config Controller) {
// 	c.Controller = config
// 	c.Data = make(map[string]interface{})
// 	//c.Data["Body"] = ""
// 	typ := c.Cntr.Type()
// 	if c.Denied {
// 		return
// 	}
// 	if method, exists := typ.MethodByName("BeforeFunc"); exists {
// 		c.Denied = method.Func.Call([]reflect.Value{c.Cntr})[0].Interface().(bool)
// 	}
// 	if c.Denied {
// 		return
// 	}
// 	c.Cntr.MethodByName(c.Action).Call([]reflect.Value{})

// 	if c.Denied {
// 		return
// 	}
// 	c.Cntr.MethodByName("Close").Call([]reflect.Value{})
// 	if c.Denied {
// 		return
// 	}
// 	if method, exists := typ.MethodByName("AfterFunc"); exists {
// 		method.Func.Call([]reflect.Value{c.Cntr})
// 	}

// }
// func (c *BaseController) Close() {
// 	if !c.Rendered {
// 		c.Render()
// 	}
// }

// // func (c *BaseController) Init(config Controller) {
// // 	c.Controller = config
// // 	c.Data = make(map[string]interface{})
// // }
// // func (c *BaseController) End() {
// // 	if !c.Rendered {
// // 		c.AutoRender()
// // 	}
// // }

// func (c *BaseController) executeTemplate() {
// 	err := c.App.Templates.ExecuteTemplate(*c.W, c.Layout, c.Data)
// 	if err != nil {
// 		fmt.Fprintf(*c.W, "%s", err)
// 	}
// }
// func (c *BaseController) setLayout() {
// 	switch {
// 	case len(c.Layout) == 0:
// 		lv := "/layouts/" + strings.Title(c.Name) + ".go." + c.Format
// 		if t := c.App.Templates.Lookup(lv); t != nil {
// 			c.Layout = lv
// 		} else {
// 			c.Layout = "/layouts/Application.go." + c.Format
// 		}
// 	case !path.IsAbs(c.Layout):
// 		c.Layout = "/layouts/" + c.Layout + ".go." + c.Format

// 	}
// 	//	c.Data["Layout"] = c.Layout
// }
// func (c *BaseController) setView() {
// 	switch {
// 	case len(c.View) == 0:
// 		c.View = "/" + c.Name + "/" + c.Action + ".go." + c.Format
// 	case !path.IsAbs(c.View):
// 		c.View = "/" + c.Name + "/" + c.View + ".go." + c.Format
// 	}
// 	//c.Data["View"] = c.View
// }

// func (c *BaseController) Render(args ...interface{}) {
// 	fmt.Println(c.Format)
// 	switch c.Format {
// 	case "json":
// 		c.RenderJson(args...)
// 	case "html":
// 		c.RenderHtml(args...)
// 	default:

// 	}
// 	c.Rendered = true
// }
// func (c *BaseController) RenderJson(args ...interface{}) {
// 	if len(args) > 0 {
// 		if b, err := json.Marshal(args); err == nil {

// 			fmt.Fprintf(*c.W, "%s", b)
// 		} else {
// 			fmt.Fprintf(*c.W, "%s", err.Error())
// 		}
// 		c.Rendered = true
// 		return
// 	}
// 	if len(c.Data) > 0 {
// 		if b, err := json.Marshal(c.Data); err == nil {
// 			fmt.Fprintf(*c.W, "%s", string(b))
// 		} else {
// 			fmt.Fprintf(*c.W, "%s", err.Error())
// 		}
// 		c.Rendered = true
// 	}
// }
// func (c *BaseController) RenderHtml(args ...interface{}) {
// 	c.setLayout()
// 	c.setView()
// 	var template_name string
// 	if c.SkipLayout {
// 		template_name = c.View
// 	} else {
// 		template_name = c.Layout
// 		c.Data["View"] = c.View
// 	}

// 	if len(args) > 0 {
// 		err := c.App.Templates.ExecuteTemplate(*c.W, template_name, args)
// 		if err != nil {
// 			fmt.Fprintf(*c.W, "%s", err)
// 		}
// 		c.Rendered = true
// 		return
// 	}
// 	if len(c.Data) > 0 {
// 		c.Data["Title"] = c.Name
// 		err := c.App.Templates.ExecuteTemplate(*c.W, template_name, c.Data)
// 		if err != nil {
// 			fmt.Fprintf(*c.W, "%s", err)
// 		}
// 		c.Rendered = true
// 	}
// }

// // func (c *BaseController) InternalRedirectTo() {
// // 	c.Rendered = true
// // 	c.AutoRender()
// // }
// func (c *BaseController) RedirectTo(url string, params map[string]interface{}) {
// 	http.Redirect(*c.W, c.R, url, http.StatusMovedPermanently)
// 	c.Denied = true
// 	//c.AutoRender()
// }
