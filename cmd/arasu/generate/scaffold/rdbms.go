package scaffold

import (
	"bytes"
	"fmt"
	"github.com/arasuresearch/arasu/lib"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (c *Scaffold) RdbmsRun() error {
	c.Attrs, c.ClientAttrs, c.ClientModelViewAttrs, c.ClientModelMetadata = parseRdbmsAttrs(c.ParseArgs)

	files := map[string]string{}
	dir := path.Join(c.SkeletonDir, "rdbms/server")

	err := filepath.Walk(dir, func(src string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		_, fn := filepath.Split(src)
		ext := filepath.Ext(fn)
		fn = strings.TrimSuffix(fn, ext)
		fn, _ = lib.ParseAndExecuteTemplateText(fn, c)
		fns := strings.Split(fn, ".")
		fn = strings.Join(fns[1:], "/") + ext
		dst := path.Join(c.AppSrcDir, fn)
		exists, _ := lib.IsExist(dst)
		if exists && fns[0] == "once" {
			return nil
		}
		if exists && !c.Force {
			return fmt.Errorf("%s already available if you want to overwrite it you can pass --force option", dst)
		}
		if ext == ".link" {
			byt, _ := ioutil.ReadFile(src)
			src = path.Join(c.App.ArasuRoot, string(byt))
		}
		files[src] = dst
		return nil
	})
	if err != nil {
		return err
	}

	for src, dst := range files {
		if err := lib.CreateTemplatedFile(src, dst, c); err != nil {
			return err
		}
		fmt.Println("created ", dst)
	}
	schemaStructureFile := path.Join(c.AppSrcDir, "server/dstores/rdbms/schema_structures.go")
	if data, err := ioutil.ReadFile(schemaStructureFile); err != nil {
		return err
	} else {
		td := []byte("\ntype " + c.Cname + " struct {Id int64}\n")
		if err = ioutil.WriteFile(schemaStructureFile, append(data, td...), os.ModePerm); err != nil {
			return err
		}
	}

	// lib.ParseImports(s)
	// 	_ "ds/bigdata/migrate"
	// _ "ds/rdbms/migrate"
	migDir := path.Join(c.App.Root, "src/tmp/arasu/main.go")
	_ = lib.AddImports(migDir, `_ "ds/rdbms/migrate"`)

	return c.RdbmsCopyClient()
	// if err := c.App.Cmd("arasu dstore rdbms migrate").Run(); err != nil {
	// 	c.revert(files)
	// 	fmt.Println("on migrating tables", err)
	// 	return
	// }
	// if err := c.App.Cmd("arasu update schema --dstore rdbms").Run(); err != nil {
	// 	c.revert(files)
	// 	fmt.Println("on updating schema structute", err)
	// 	return
	// }

	//rm -rf d0;arasubuild;arasu new d0 -d mysql -ds rdbms;arasu generate scaffold Post name
	//arasu generate scaffold Post name

}

func (c *Scaffold) RdbmsCopyClient() error {
	var formSrc, formDst string
	files := map[string]string{}
	err := filepath.Walk(path.Join(c.SkeletonDir, "rdbms/client"), func(src string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		_, fn := filepath.Split(src)
		ext := filepath.Ext(fn)
		fn = strings.TrimSuffix(fn, ext)
		fn, _ = lib.ParseAndExecuteTemplateText(fn, c)
		fns := strings.Split(fn, ".")
		fn = strings.Join(fns, "/") + ext
		dst := path.Join(c.AppSrcDir, "client", fn)
		if ext == ".link" {
			byt, _ := ioutil.ReadFile(src)
			src = path.Join(c.App.ArasuRoot, string(byt))
		}
		if strings.HasSuffix(src, "form.html") {
			formSrc, formDst = src, dst
			return nil
		}
		files[src] = dst
		return nil
	})
	if err != nil {
		return err
	}
	for src, dst := range files {
		if err := lib.CreateTemplatedFile(src, dst, c); err != nil {
			return err
		}
		fmt.Println("created ", dst)
	}
	return c.RdbmsCopyClientForm(formSrc, formDst)
}
func (c *Scaffold) RdbmsCopyClientForm(src, dst string) error {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	Templates, err := template.New("tmp").Parse(string(content))
	if err != nil {
		return err
	}
	if err = loadClientViewIndividualTemplates(Templates, path.Join(c.SkeletonDir, "common")); err != nil {
		return err
	}
	if err = Templates.Execute(&buf, c); err != nil {
		return err
	}
	if err := lib.CreateAndWriteFile(dst, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func loadClientViewIndividualTemplates(templates *template.Template, dir string) error {
	var filenames []string
	err := filepath.Walk(dir, func(src string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filenames = append(filenames, src)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		_, fn := filepath.Split(filename)
		fname := strings.TrimSuffix(fn, filepath.Ext(fn))
		tl := templates.New(fname)
		if _, err = tl.Parse(string(b)); err != nil {
			return err
		}
	}
	template.Must(templates, nil)
	return nil
}
