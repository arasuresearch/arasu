package migrate

import (
	. "github.com/arasuresearch/arasu/datastorage/rdbms"
	"github.com/arasuresearch/arasu/datastorage/rdbms/adapter/abstract"
)

func init() {
	Add("M{{.MigrationVersion}}_Create_{{.Name}}s", Migration{
		Up: func() {
			CreateTable("{{.Name}}s", func(t *abstract.Table) {
				{{range $k,$v := .Attrs}}t.{{$v}}("{{$k}}")
				{{end}}
				t.Timestamps()
			})
			AddIndex("{{.Name}}", []string{"id"})
		},

		Down: func() {
			DropTable("{{.Name}}")
			DropIndex("{{.Name}}", []string{"id"})
		},
	})
}

