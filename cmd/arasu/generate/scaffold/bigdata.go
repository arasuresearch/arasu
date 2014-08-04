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

func (c *Scaffold) BigdataRun() error {
	c.Attrs, c.ClientAttrs, c.ClientModelViewAttrs, c.ClientModelMetadata = parseBigDataAttrs(c.ParseArgs)
	//fmt.Println(c.Attrs, c.ClientAttrs, c.ClientModelViewAttrs, c.ClientModelMetadata)
	files := map[string]string{}
	dir := path.Join(c.SkeletonDir, "bigdata/server")
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
	schemaStructureFile := path.Join(c.AppSrcDir, "server/dstores/bigdata/schema_structures.go")
	if data, err := ioutil.ReadFile(schemaStructureFile); err != nil {
		return err
	} else {
		td := []byte("\ntype " + c.Cname + " struct {Id string}\n")
		if err = ioutil.WriteFile(schemaStructureFile, append(data, td...), os.ModePerm); err != nil {
			return err
		}
	}
	//pkgName := "ds/bigdata/migrate"
	migDir := path.Join(c.App.Root, "src/tmp/arasu/main.go")
	_ = lib.AddImports(migDir, `_ "ds/bigdata/migrate"`)
	return c.BigdataCopyClient()

	//fmt.Println("==================================")

	//rm -rf d0;arasubuild;arasu new d0 -d mysql -ds rdbms;arasu generate scaffold Post name
	//arasu generate scaffold Post name

}
func (c *Scaffold) revert(files map[string]string) {
	for _, v := range files {
		_ = os.Remove(v)
	}
}
func (c *Scaffold) BigdataCopyClient() error {
	var formSrc, formDst string
	files := map[string]string{}
	err := filepath.Walk(path.Join(c.SkeletonDir, "bigdata/client"), func(src string, info os.FileInfo, err error) error {
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
	return c.BigdataCopyClientForm(formSrc, formDst)
}
func (c *Scaffold) BigdataCopyClientForm(src, dst string) error {
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
	//fmt.Println(c.ClientModelViewAttrs)
	if err = Templates.Execute(&buf, c); err != nil {
		return err
	}
	if err := lib.CreateAndWriteFile(dst, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

//git clean -f -d;arasubuild;arasu generate scaffold Mike Contact Post:{Fn,Ln} Profile:{Name:String,Age:int,Dob:DateTime}
// arasubuild;arasu update schema --ds-path bigdata
// arasu update schema --ds-path bigdata
//  rm -rf d1; arasubuild;arasu new d1 -d hbase -ds bigdata
