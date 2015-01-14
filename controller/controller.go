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
