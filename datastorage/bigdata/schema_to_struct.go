package bigdata

import (
	"bytes"
	"github.com/arasuresearch/arasu/datastorage/bigdata/adapter"
	"github.com/arasuresearch/arasu/lib"
	"os/exec"
	"path"
	"text/template"
)

func SchemaToStruct(name, conf, modelDir, templatePath string) error {
	adap := adapter.New(name, conf)
	STS, err := adap.SchemaToStruct(modelDir)
	if err != nil {
		return err
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	var w bytes.Buffer
	if err = t.Execute(&w, STS); err != nil {
		return err
	}
	fout := path.Join(modelDir, "schema_structures.go")
	if err := lib.CreateAndWriteFile(fout, w.Bytes()); err != nil {
		return err
	}
	_, _ = exec.Command("go", "fmt", fout).CombinedOutput()
	return nil
}
