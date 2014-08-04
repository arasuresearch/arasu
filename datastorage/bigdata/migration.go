package bigdata

import (
	"bytes"
	"github.com/arasuresearch/arasu/datastorage/bigdata/adapter"
	"github.com/arasuresearch/arasu/datastorage/bigdata/adapter/abstract"
	"github.com/arasuresearch/arasu/lib"
	"log"
	"os/exec"
	"path"
	"text/template"
)

// migration type to a specific version

type Migration struct {
	Up   func()
	Down func()
}

var DC *DstoreCommand

//list of migrations with verions
var Migrations = map[string]Migration{}

func Add(version string, migration Migration) {
	if _, ok := Migrations[version]; ok {
		log.Fatal("already same version migration are defined")
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
	// if !d.Adapter.IsDatabaseExists() {
	// 	log.Fatalf("'%s' namespace not available , please create it first", d.Adapter.GetDbName())
	// }
	versions, err := d.Adapter.GetAllSchemaMigration()
	if err != nil {
		return err
	}
	//fmt.Println(versions)
	for k, v := range d.Migrations {
		if !lib.StringArrayContains(versions, k) {
			v.Up()
			d.Adapter.InsertIntoSchemaMigration(k)
		}
	}
	d.Adapter.DumpSchema()
	return nil
}

func (d *DstoreCommand) Rollback() error {
	// if !d.Adapter.IsDatabaseExists() {
	// 	log.Fatalf("'%s' namespace not available , please create it first", d.Adapter.GetDbName())
	// }
	versions, err := d.Adapter.GetAllSchemaMigration()
	if err != nil {
		return err
	}

	for k, v := range d.Migrations {
		if lib.StringArrayContains(versions, k) {
			v.Down()
			d.Adapter.DeleteFromSchemaMigration(k)
		}
	}
	d.Adapter.DumpSchema()
	return nil
}

func CreateTable(name string, callback func(t *abstract.Table)) error {
	return DC.Adapter.CreateTable(name, callback)
}
func DropTable(name string) error {
	return DC.Adapter.DropTable(name)
}
