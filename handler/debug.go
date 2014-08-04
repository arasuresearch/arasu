package handler

import (
	"encoding/json"
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"github.com/arasuresearch/arasu/controller"
	"github.com/arasuresearch/arasu/router"
	"net/http"
	"reflect"
)

type DebugHandler struct {
	App      *app.App
	Routes   *router.Router
	Registry router.Registry
}

func (h *DebugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "index.html", 302)
		return
	}
	if err := handle(w, r, h.Registry, h.Routes, h.App); err != nil {
		fmt.Println(err.Error())
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func handle(w http.ResponseWriter, r *http.Request, registry router.Registry, routes *router.Router, App *app.App) error {
	//buep is  base_url_embeded_params
	hand, _, fname, _, buep, err := routes.Find(&w, r, registry)
	//fmt.Println(hand, cname, fname, format, buep)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "%s", err.Error())
		return err
	}
	//fmt.Println(r.Header.Get("accept"))
	//fmt.Println(hand, fname)
	//cntr := reflect.New(hand.Type)
	cntr := reflect.New(hand.Type.Elem())
	Params, err := parseParams(r, buep)
	//fmt.Println(Params, err)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return err
	}

	config := controller.Controller{
		W:      &w,
		R:      r,
		Action: fname,
		App:    App,
		Params: Params,
	}

	cntr.Elem().FieldByName("Controller").Set(reflect.ValueOf(config))
	action := cntr.MethodByName(fname)

	if before := cntr.MethodByName("BeforeFunc"); before.IsValid() {
		res := before.Call(nil)
		if len(res) > 0 && res[len(res)-1].Interface() != nil {
			err := res[len(res)-1].Interface().(error)
			fmt.Fprintf(w, "%s", err)
			return err
		}
	}
	////////////////////////

	args, err := parseArgsFromParams(action.Type(), hand.FuncsArgs[fname], Params)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return err
	}

	result := action.Call(args)

	if after := cntr.MethodByName("AfterFunc"); after.IsValid() {
		res := after.Call(nil)
		if len(res) > 0 && res[len(res)-1].Interface() != nil {
			err := res[len(res)-1].Interface().(error)
			fmt.Fprintf(w, "%s", err)
			return err
		}
	}
	datas := make([]interface{}, len(result))
	for i, e := range result {
		switch v := e.Interface().(type) {
		case error:
			datas[i] = v.Error()
		default:
			datas[i] = v
		}
	}
	if b, err := json.Marshal(datas); err == nil {
		fmt.Fprintf(w, "%s", b)
	} else {
		fmt.Fprintf(w, "%s", err)
		return err
	}
	var shower = make(map[interface{}]interface{})
	for _, e := range datas {
		v0 := fmt.Sprintf("%v", e)
		if len(v0) > 40 {
			v0 = v0[:40] + "..."
		}
		shower[fmt.Sprintf("%T", e)] = v0
	}
	// data_names := make([]string, len(datas))
	// for i, e := range datas {
	// 	data_names[i] = fmt.Sprintf("%T", e)
	// }

	fmt.Println(r.Method, r.URL.Path, hand.Type, fname, Params, "-->", shower)

	return nil
}
