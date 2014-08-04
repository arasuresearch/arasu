package abstract

import (
	"bitbucket.org/pkg/inflect"
	"fmt"
	"github.com/arasuresearch/arasu/lib"
	"io/ioutil"
	"path"
	"reflect"
	"strings"
)

type SchemaToStruct struct {
	Pkg     string
	Imports lib.AS
	Tables  []TableStruct
}

type TableStruct struct {
	Name    string
	Columns []ColumnStruct
	Tag     string
}

type ColumnStruct struct {
	Name, Type, Tag string
}

func (a *AbstractAdapter) SchemaToStruct(adapter_path string) (*SchemaToStruct, error) {

	model_names, err := ioutil.ReadDir(adapter_path)
	if err != nil {
		return nil, err
	}
	models := lib.AS{}
	for _, e := range model_names {
		if !e.IsDir() && e.Name() != "init.go" && e.Name() != "schema_structures.go" {
			models.Add(strings.TrimSuffix(e.Name(), ".go"))
		}
	}
	fmt.Println(models)
	imports := lib.AS{}
	tables := []TableStruct{}

	for _, e := range models {
		table := TableStruct{Name: strings.Title(e)}
		//fmt.Println("SHOW FIELDS FROM `" + e + "`")

		tname := a.Quote(inflect.Pluralize(e))

		rows, err := a.DB.Query("SHOW FIELDS FROM " + tname)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			args := make([]interface{}, len(columns))
			for i, _ := range columns {
				var l interface{}
				args[i] = &l
			}
			if err := rows.Scan(args...); err != nil {
				return nil, err
			}
			var name, typ, tag string
			tags := []string{}
			var primary bool
			if field := reflect.Indirect(reflect.ValueOf(args[0])).Interface(); field != nil {
				v0 := inflect.Camelize(string(field.([]byte)))
				//v0 := strings.Title(string(field.([]byte)))
				name = v0
			}
			if field := reflect.Indirect(reflect.ValueOf(args[1])).Interface(); field != nil {
				v0 := string(field.([]byte))
				switch {
				case strings.Contains(v0, "int"):
					switch v0 {
					case "int(1)", "tinyint(1)":
						typ = "bool"
					case "int(11)":
						typ = "int64"
					default:
						typ = "int"
					}
				case strings.Contains(v0, "decimal"), strings.Contains(v0, "float"):
					typ = "float32"
				case strings.Contains(v0, "bool"), strings.Contains(v0, "int(1)"):
					typ = "bool"
				case strings.Contains(v0, "varchar"):
					typ = "string"
				case strings.Contains(v0, "byte"), strings.Contains(v0, "binary"):
					typ = "[]byte"
				case strings.Contains(v0, "datetime"), strings.Contains(v0, "date"):
					typ = "time.Time"
					imports.Add("time")
				}
			}
			if field := reflect.Indirect(reflect.ValueOf(args[2])).Interface(); field != nil {
				v0 := inflect.Camelize(string(field.([]byte)))
				//v0 := strings.Title(string(field.([]byte)))
				if v0 == "NO" {
					tags = append(tags, "not null")
				}
			}
			if field := reflect.Indirect(reflect.ValueOf(args[3])).Interface(); field != nil {
				v0 := string(field.([]byte))
				if v0 == "PRI" {
					primary = true
				}
			}
			if field := reflect.Indirect(reflect.ValueOf(args[4])).Interface(); field != nil {
				v0 := string(field.([]byte))
				if len(v0) > 0 {
					tags = append(tags, "default:"+a.SingleQuote(v0))
				}
			}
			if len(tags) > 0 {
				tag = `sql:"` + strings.Join(tags, ";") + `"`
			}
			if primary {
				tag += ` primaryKey:"yes"`
			}
			if len(tag) > 0 {
				tag = a.Quote(tag)
			}

			column := ColumnStruct{name, typ, tag}
			table.Columns = append(table.Columns, column)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	_, pkg := path.Split(adapter_path)
	//pkg = strings.Title(pkg)

	data := SchemaToStruct{pkg, imports, tables}
	return &data, nil
}

// func write(data Adapter) error {
// 	t, err := template.ParseFiles(path.Join(Conf.ArasuRoot, "skeleton/templates/schema.go.tmpl"))
// 	if err != nil {
// 		return err
// 	}

// 	var w bytes.Buffer
// 	if err := t.Execute(&w, data); err != nil {
// 		return err
// 	}

// 	fout := path.Join(rdbmsPath, data.Pkg, "schema_structures.go")

// 	if err := lib.CreateAndWriteFile(fout, w.Bytes()); err != nil {
// 		return err
// 	}
// 	_, _ = exec.Command("go", "fmt", fout).CombinedOutput()
// 	return nil
// }
