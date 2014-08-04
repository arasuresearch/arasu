package migrate

import (
	. "github.com/arasuresearch/arasu/datastorage/bigdata"
	"github.com/arasuresearch/arasu/datastorage/bigdata/adapter/abstract"
)

func init() {
	Add("M{{.MigrationVersion}}_Create_{{.Cname}}", Migration{
		Up: func() {
			CreateTable("{{.Cname}}", func(t *abstract.Table) {
				{{range $k,$v := .Attrs}}t.{{$v}}("{{$k}}")
				{{end}}
			})
		},
		Down: func() {
			DropTable("{{.Cname}}")
		},
	})
}
