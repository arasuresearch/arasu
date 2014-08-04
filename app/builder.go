package app

import (
	"bytes"
	"fmt"
	"github.com/arasuresearch/arasu/lib"
	"github.com/arasuresearch/arasu/watcher"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

type Builder struct {
	App          *App
	CntrTemplate *template.Template
	CntrRoot     string
	CntrPkgs     []string
	Watcher      *watcher.Watcher
	DataCntrs    map[string]*DataCntr
	DataFiles    map[string][]byte
	Restart      bool
}

func (a *App) NewBuilder() (*Builder, error) {
	b := &Builder{App: a}
	return b, b.init()
}
func (b *Builder) init() error {
	b.CntrTemplate, _ = template.ParseFiles(path.Join(b.App.ArasuRoot, "skeleton/templates/cntr.go.tmpl"))
	b.CntrRoot = path.Join(b.App.Root, "src/server/controllers")
	b.Watcher, _ = watcher.New(path.Join(b.App.Root, "src/server"))
	//b.DataFiles = make(map[string][]byte)
	b.DataCntrs = make(map[string]*DataCntr)
	cntrs := lib.AS{}
	err := filepath.Walk(b.Watcher.Dir, func(src string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(src) == ".go" {
			rsrc := strings.TrimPrefix(src, b.Watcher.Dir)
			//b.DataFiles[rsrc], _ = ioutil.ReadFile(src)
			if strings.HasPrefix(rsrc, "/controllers") {
				url, _ := filepath.Split(strings.TrimPrefix(src, b.CntrRoot))
				url = filepath.Clean(url)
				if cntrs.Add(url) {
					cntr_pkg_root, _ := filepath.Split(src)
					b.DataCntrs[url] = create_data_cntr(url, cntr_pkg_root)
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, v := range b.DataCntrs {
		b.create_dispatcher(v)
	}
	return b.App.Cmd("go install tmp/arasu/debug").Run()
}

func (b *Builder) create_dispatcher(dc *DataCntr) {
	var buf bytes.Buffer
	_ = b.CntrTemplate.Execute(&buf, dc)
	rdst := strings.Join(strings.Split(filepath.Clean(dc.Url[1:]), "/"), ".") + ".go"
	dst := path.Join(b.App.Root, "src/tmp/dispatchers", rdst)
	_ = lib.CreateAndWriteFile(dst, buf.Bytes())
}

func (b *Builder) dispatch() {
	for _, e := range b.Watcher.Pkgs {
		if strings.HasPrefix(e, b.CntrRoot) {

			url := filepath.Clean(strings.TrimPrefix(e, b.CntrRoot))
			if b.remove_dispatchers_for_deleted_controllers(url, e) {
				continue
			}

			dataCntr := create_data_cntr(url, e)
			local := false
			if dc, ok := b.DataCntrs[url]; ok {
				if !reflect.DeepEqual(dc, dataCntr) {
					local = true
				}
			} else {
				local = true
			}
			if local {
				b.create_dispatcher(dataCntr)
				b.DataCntrs[url] = dataCntr
				b.Restart = true
			}
		}
	}
}
func (b *Builder) remove_dispatchers_for_deleted_controllers(key, cntr_root string) bool {
	err := filepath.Walk(cntr_root, func(src string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if path.Ext(src) == ".go" {
			return fmt.Errorf("go files are available")
		}
		return nil
	})
	if err != nil {
		return false
	}
	disFile := path.Join(b.App.Root, "src/tmp/dispatchers", path.Base(cntr_root)+".go")
	_ = os.Remove(disFile)
	delete(b.DataCntrs, key)
	return true
}

func (b *Builder) ReBuild() error {
	b.dispatch()
	b.Watcher.Clean()
	return b.App.Cmd("go install tmp/arasu/debug").Run()
}
