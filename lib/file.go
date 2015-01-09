package lib

import (
	"bytes"
	"fmt"
	"github.com/arasuresearch/arasu/lib/stringer"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

func CopyFile(src, dst string) (int64, error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()
	dir, _ := path.Split(dst)
	os.MkdirAll(dir, os.ModePerm)

	df, err0 := os.Create(dst)
	if err0 != nil {
		return 0, err0
	}
	defer df.Close()
	return io.Copy(df, sf)
}

func CreateAndWriteFile(name string, data []byte) error {
	dir, _ := path.Split(name)
	os.MkdirAll(dir, os.ModePerm)
	return ioutil.WriteFile(name, data, os.ModePerm)
}

func CreateTemplatedFile(src, dst string, data interface{}) error {
	var TempWriter bytes.Buffer
	fileTemplate, err := template.ParseFiles(src)
	if path.Ext(src) == ".tmpl" {
		dst = strings.TrimSuffix(dst, ".tmpl")
		content, err := ioutil.ReadFile(src)
		if err != nil {
			return err
		}
		parsedContent, err := ParseAndExecuteTemplateText(string(content), data)
		if err != nil {
			return err
		}
		parsedContent = strings.Replace(parsedContent, "_ocb_", "{{", -1)
		parsedContent = strings.Replace(parsedContent, "_ccb_", "}}", -1)
		err = CreateAndWriteFile(dst, []byte(parsedContent))
		if err != nil {
			return err
		}
		return nil
	}

	// if strings.HasSuffix(path.Ext(src), "2") {
	// 	dst = strings.TrimSuffix(dst, "2")
	// 	content, err := ioutil.ReadFile(src)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	parsedContent, err := ParseAndExecuteTemplateText(string(content), data)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	parsedContent = strings.Replace(parsedContent, "_ocb_", "{{", -1)
	// 	parsedContent = strings.Replace(parsedContent, "_ccb_", "}}", -1)
	// 	err = CreateAndWriteFile(dst, []byte(parsedContent))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }

	if err != nil {

		// parsedContent, err := ParseDoubleTemplateString(string(content), data)
		// if err != nil {
		// 	return err
		// }
		// err = ioutil.WriteFile(dst, []byte(parsedContent), os.ModePerm)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println("Error occured while parsing ", src, " to ", dst, "---->", err)
		// fmt.Println("So directly copying the files")
		// if strings.Contains(err.Error(), "unexpected bad number syntax") {
		// 	_, e0 := CopyFile(src, dst)
		// 	return e0
		// } else if path.Ext(src) == ".html" &&
		// 	strings.Contains(err.Error(), "function") &&
		// 	strings.Contains(err.Error(), "not defined") {
		// 	//fmt.Println(src, dst, "==========", err)
		// 	_, e0 := CopyFile(src, dst)
		// 	return e0
		// } else {
		// 	return err
		// }
		//fmt.Println("----", err.Error(), "--------")
	}

	if err := fileTemplate.Execute(&TempWriter, data); err != nil {
		fmt.Println(err)
		return err
	}
	if err := CreateAndWriteFile(dst, TempWriter.Bytes()); err != nil {
		return err
	}
	return nil
}

func ParseAndExecuteTemplateText(text string, data interface{}) (string, error) {
	var TempWriter bytes.Buffer

	textTemplate, err := template.New("tmp").Parse(text)
	if err != nil {
		return "", err
	}
	err = textTemplate.Execute(&TempWriter, data)
	if err != nil {
		return "", err
	}
	return TempWriter.String(), nil
}

func IsExist(src string) (bool, error) {
	if !path.IsAbs(src) {
		pwd, _ := os.Getwd()
		src = path.Join(pwd, src)
	}
	_, err := os.Stat(src)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func AddImports(src, imp string) error {
	if imps, err := ParseImports(src); err != nil {
		return err
	} else if stringer.Contains(imps, imp) {
		return nil
	}
	data, _ := ioutil.ReadFile(src)
	rows := strings.Split(string(data), "\n")
	res := []string{}
	index := 0
	for i, e := range rows {
		if strings.Contains(e, "import") {
			index = i + 1
			break
		}
	}
	res = append(res, rows[:index]...)
	res = append(res, imp)
	res = append(res, rows[index:]...)
	CreateAndWriteFile(src, []byte(strings.Join(res, "\n")))
	if _, err := exec.Command("go", "fmt", src).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

func ParseImports(s string) ([]string, error) {
	fs := token.NewFileSet()
	imports := []string{}
	var importName string
	switch path.Ext(s) {
	case ".go":
		f, err := parser.ParseFile(fs, s, nil, parser.ImportsOnly)
		if err != nil {
			return imports, err
		}
		for _, e := range f.Imports {
			if e.Name != nil {
				importName = e.Name.String() + ` ` + e.Path.Value
			} else {
				importName = e.Path.Value
			}
			imports = append(imports, importName)
		}
	case "":
		fun := func(f os.FileInfo) bool {
			return !f.IsDir() &&
				!strings.HasPrefix(f.Name(), ".") &&
				path.Ext(f.Name()) == ".go"
		}
		pkgs, err := parser.ParseDir(fs, s, fun, parser.ImportsOnly)
		if err != nil {
			return imports, err
		}
		for _, v := range pkgs {
			for _, v0 := range v.Files {
				for _, e := range v0.Imports {
					if e.Name != nil {
						importName = e.Name.String() + ` ` + e.Path.Value
					} else {
						importName = e.Path.Value
					}
					imports = append(imports, importName)
				}
			}
		}
	default:

	}
	return imports, nil
}
