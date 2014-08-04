package rdbms

import (
	"bytes"
	"fmt"
	"github.com/arasuresearch/arasu/datastorage/rdbms/adapter"
	"github.com/arasuresearch/arasu/lib"
	"log"
	"os/exec"
	"path"
	"text/template"
)

type Migration struct {
	Up   func()
	Down func()
}

var DC *DstoreCommand

//list of migrations with verions
var Migrations = map[string]Migration{}

func Add(version string, migration Migration) {
	if _, ok := Migrations[version]; ok {
		log.Fatal("already same version of migration are defined")
	}
	Migrations[version] = migration
}

// migrator command
type DstoreCommand struct {
	AdapterName        string
	Conf               string
	Migrations         map[string]Migration
	Cmd                string
	Args               string
	Adapter            adapter.Adapter
	ModelRoot          string
	SchemaTemplatePath string
}

//starting migration like
// arasu ds bd [create,drop,migrate,rollback,etc]
func (d *DstoreCommand) Start() (err error) {
	d.Migrations = Migrations
	d.Adapter = adapter.New(d.AdapterName, d.Conf)
	DC = d
	switch d.Args {
	case "create":
		err = d.Adapter.CreateDatabase()
	case "drop":
		err = d.Adapter.DropDatabase()
	default:
		if err = d.Adapter.Init(); err != nil {
			return
		}
		switch d.Args {
		case "migrate":
			err = d.Migrate()
		case "rollback":
			err = d.Rollback()
		case "migrate:up":

		case "migrate:down":
		case "seed":
		default:
		}
		if err == nil {
			err = d.SchemaToStruct()
		}
	}
	return
}
func (d *DstoreCommand) SchemaToStruct() error {
	STS, err := d.Adapter.SchemaToStruct(d.ModelRoot)
	if err != nil {
		return err
	}
	t, err := template.ParseFiles(d.SchemaTemplatePath)
	if err != nil {
		return err
	}
	var w bytes.Buffer
	if err = t.Execute(&w, STS); err != nil {
		return err
	}
	fout := path.Join(d.ModelRoot, "schema_structures.go")
	if err := lib.CreateAndWriteFile(fout, w.Bytes()); err != nil {
		return err
	}
	_, _ = exec.Command("go", "fmt", fout).CombinedOutput()
	return nil
}

func (d *DstoreCommand) Migrate() error {
	versions, err := d.Adapter.GetAllSchemaMigration()
	if err != nil {
		return err
	}
	for k, v := range d.Migrations {
		if !lib.StringArrayContains(versions, k) {
			if err := d.Adapter.Transaction(v.Up); err != nil {
				return err
			}
			if err := d.Adapter.InsertIntoSchemaMigration(k); err != nil {
				return err
			}
			fmt.Printf("%s migrated\n", k)

		}
	}
	d.Adapter.DumpSchema()
	return nil
}

func (d *DstoreCommand) Rollback() error {
	versions, err := d.Adapter.GetAllSchemaMigration()
	//fmt.Println(versions, err)
	if err != nil {
		return err
	}
	for k, v := range d.Migrations {
		if lib.StringArrayContains(versions, k) {
			if err := d.Adapter.Transaction(v.Down); err != nil {
				return err
			}
			if err := d.Adapter.DeleteFromSchemaMigration(k); err != nil {
				return err
			}
			fmt.Printf("%s reverse migrated \n", k)
		}
	}
	d.Adapter.DumpSchema()
	return nil
}
func CreateTable(name string, args ...interface{}) {
	DC.Adapter.CreateTable(name, args...)
}
func AlterTable(name string, args ...interface{}) {
	DC.Adapter.AlterTable(name, args...)
}
func DropTable(name string, args ...interface{}) {
	DC.Adapter.DropTable(name, args...)
}

func AddIndex(table_name string, column_names []string, args ...interface{}) {
	var options map[string]interface{}
	if len(args) > 1 {
		options = args[0].(map[string]interface{})
	}
	DC.Adapter.CreateIndex(table_name, column_names, options)
}
func AlterIndex(name string, args ...interface{}) {
	DC.Adapter.AlterTable(name, args...)
}

func DropIndex(table_name string, column_names []string, args ...interface{}) {
	var options map[string]interface{}
	if len(args) > 1 {
		options = args[0].(map[string]interface{})
	}
	DC.Adapter.DropIndex(table_name, column_names, options)
}
