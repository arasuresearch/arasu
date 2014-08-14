package view

import (
	"github.com/arasuresearch/arasu/lib"
	"github.com/arasuresearch/arasu/lib/stringer"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ParseViews(dir string, allowed []string) (*template.Template, error) {
	name := filepath.Base(dir)
	templates := template.New(name)
	var filenames []string
	view_path := filepath.Join(dir, "src/server/views")
	allowed = lib.AS{".html", ".json", ".xml"}

	err := filepath.Walk(view_path, func(src string, info os.FileInfo, err error) error {
		if !info.IsDir() && stringer.Contains(allowed, filepath.Ext(src)) {
			filenames = append(filenames, src)
		}
		return nil
	})
	if err != nil {
		return templates, err
	}
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return templates, err
		}
		src := strings.TrimPrefix(filename, view_path)
		name := filepath.Clean(src)
		tl := templates.New(name)
		_, err = tl.Parse(string(b))
		if err != nil {
			return templates, err
		}
	}
	return template.Must(templates, nil), nil
}
