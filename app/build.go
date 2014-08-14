package app

import (
	"bytes"
	"github.com/arasuresearch/arasu/lib"
	"github.com/arasuresearch/arasu/lib/stringer"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type DataCntr struct {
	PkgPath   string
	Cntr      string
	Url       string
	FuncsArgs map[string][]string
}

func (a *App) Build() error {
	cntr_template, err := template.ParseFiles(path.Join(a.ArasuRoot, "skeleton/templates/cntr.go.tmpl"))
	if err != nil {
		return err
	}
	cntr_root := path.Join(a.Root, "src/server/controllers")
	data_cntrs := make(map[string]*DataCntr)
	cntrs := lib.AS{}
	err = filepath.Walk(cntr_root, func(src string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(src) == ".go" {
			url, _ := filepath.Split(strings.TrimPrefix(src, cntr_root))
			url = filepath.Clean(url)
			if cntrs.Add(url) {
				cntr_pkg_root, _ := filepath.Split(src)
				data_cntrs[url] = create_data_cntr(url, cntr_pkg_root)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	dsts := []string{}
	for _, v := range data_cntrs {
		buf.Reset()
		if err := cntr_template.Execute(&buf, v); err != nil {
			return err
		}
		rdst := strings.Join(strings.Split(filepath.Clean(v.Url[1:]), "/"), ".") + ".go"
		dst := path.Join(a.Root, "src/tmp/dispatchers", rdst)
		if err := lib.CreateAndWriteFile(dst, buf.Bytes()); err != nil {
			return err
		}
		dsts = append(dsts, dst)
	}
	if len(dsts) > 0 {
		if err := remove_anonymous_dispatchers(path.Join(a.Root, "src/tmp/dispatchers"), dsts); err != nil {
			return err
		}
	}
	return a.Cmd("go install tmp/arasu").Run()
}

func remove_anonymous_dispatchers(disRoot string, dsts []string) error {
	return filepath.Walk(disRoot, func(src string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !stringer.Contains(dsts, src) {
			if err := os.Remove(src); err != nil {
				return err
			}
		}
		return nil
	})
}
func create_data_cntr(url, cntr_pkg_root string) *DataCntr {
	_, funcs_args, _ := parseController(cntr_pkg_root)
	pkgName := filepath.Base(url)
	cntrName := pkgName + "." + stringer.Camelize(pkgName)
	return &DataCntr{"server/controllers" + url, cntrName, url, funcs_args}
}
func parseController(cntr_path string) (string, map[string][]string, error) {
	fs := token.NewFileSet()
	filter := func(f os.FileInfo) bool {
		return !f.IsDir() && !strings.HasPrefix(f.Name(), ".") && strings.HasSuffix(f.Name(), ".go")
	}
	pkgs, err := parser.ParseDir(fs, cntr_path, filter, parser.AllErrors)
	if err != nil {
		return "", nil, err
	}
	Cfpl := make(map[string][]string)
	var cname string
	for k, v := range pkgs {
		cname = strings.Title(k) + "Controller"
		for _, v0 := range v.Files {
			for _, e := range v0.Decls {
				if fd, ok := e.(*ast.FuncDecl); ok && fd.Recv != nil && fd.Name.IsExported() {
					rt := fd.Recv.List[0].Type
					var cn string
					if t, ok := rt.(*ast.StarExpr); ok {
						cn = t.X.(*ast.Ident).Name
					} else {
						cn = rt.(*ast.Ident).Name
					}
					if cname == cn {
						fn := fd.Name.Name
						ps := fd.Type.Params.List
						Cfpl[fn] = make([]string, len(ps))
						for i, p := range ps {
							Cfpl[fn][i] = p.Names[0].String()
						}
					}
				}
			}
		}
	}
	return cname, Cfpl, nil
}
