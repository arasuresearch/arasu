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
		// if ext == ".link" {
		// 	byt, _ := ioutil.ReadFile(src)
		// 	src = path.Join(c.App.ArasuRoot, string(byt))
		// }
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

	migDir := path.Join(c.App.Root, "src/tmp/arasu/main.go")
	_ = lib.AddImports(migDir, `_ "ds/rdbms/migrate"`)
	if err = c.CopyClient(); err != nil {
		return err
	}

	return c.CopyClientIndividual("rdbms/client")

}

func (c *Scaffold) CopyClient() error {
	files := map[string]string{}
	err := filepath.Walk(path.Join(c.SkeletonDir, "common/client"), func(src string, info os.FileInfo, err error) error {
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
		if ext == ".tmpl" {
			dstl := []uint8(dst)
			dstl[strings.LastIndex(dst, "/")] = uint8('.')
			dst = string(dstl)
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
	return c.AppendBindings()
}

func (c *Scaffold) CopyClientIndividual(individual_src string) error {
	var formSrc, formDst string
	files := map[string]string{}
	err := filepath.Walk(path.Join(c.SkeletonDir, individual_src), func(src string, info os.FileInfo, err error) error {
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
		if ext == ".tmpl" {
			dstl := []uint8(dst)
			dstl[strings.LastIndex(dst, "/")] = uint8('.')
			dst = string(dstl)
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
	return c.CopyClientForm(formSrc, formDst)
}

func (c *Scaffold) CopyClientForm(src, dst string) error {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	Templates, err := template.New("tmp").Parse(string(content))
	if err != nil {
		return err
	}
	if err = loadClientViewIndividualTemplates(Templates, path.Join(c.SkeletonDir, "inputs")); err != nil {
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

func (c *Scaffold) AppendBindings() error {
	filenames := []string{
		path.Join(c.AppSrcDir, "client/lib/controllers/controllers.dart"),
		path.Join(c.AppSrcDir, "client/lib/models/models.dart"),
		path.Join(c.AppSrcDir, "client/lib/routes/routes.dart"),
	}
	for _, e := range filenames {
		if content, err := ioutil.ReadFile(e); err != nil {
			return err
		} else {
			data := Inserter(string(content), c.Name+"s")
			err = lib.CreateAndWriteFile(e, []byte(data))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Inserter(data string, newstr string) string {
	paternstr, oldstr := "//::pattern::", "+++replace_this+++"
	var result []string
	a := strings.Split(data, "\n")
	for _, e := range a {
		if strings.HasPrefix(e, paternstr) {
			s := strings.TrimPrefix(e, paternstr)
			s = strings.Replace(s, oldstr, newstr, -1)
			result = append(result, s)
		}
		result = append(result, e)
	}
	return strings.Join(result, "\n")
}

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
