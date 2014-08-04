package abstract

import (
	"github.com/arasuresearch/arasu/lib"
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
