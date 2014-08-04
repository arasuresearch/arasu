package abstract

import (
	"log"
	"strconv"
)

type Table struct {
	Name      string
	Columns   []Column
	Indexes   []Index
	Native    map[string]interface{}
	Temporary bool
	Options   map[string]interface{}
	As        string
	AbstractAdapter  *AbstractAdapter
}
type AlterTable struct {
	Name      string
	Adds      []Column
	Columns   []Column
	Indexes   []Index
	Native    bool
	Temporary bool
	Options   map[string]string
}
type DropTable struct {
	Name string
}

type Column struct {
	Name          string
	Type          string
	As            string
	Primary       bool
	Limit         uint64
	Array         bool
	Precision     int
	Scale         int
	Default       string
	NotNull       bool
	First         bool
	After         bool
	AutoIncrement bool
	Options       map[string]interface{}
}
type AlterColumn struct {
	Column  Column
	Ctype   string
	Options map[string]interface{}
}
type DropColumn struct {
	Name string
}

type Index struct {
	Table       string
	Name        string
	Unique      bool
	ColumnNames []string
	Lengths     int
	orders      []string
	where       string
	Type        string
	Using       string
}

type AlterIndex struct {
	Column  Column
	Ctype   string
	Options map[string]interface{}
}
type DropIndex struct {
	Name string
}

func (c *Column) TypeToSql(a *AbstractAdapter) string {
	if len(c.As) > 0 {
		return c.As
	} else {
		return a.TypeToSql(c.Type, c.Limit, c.Precision, c.Scale)
	}
}

func get_uint64(a interface{}) (limit uint64) {
	switch a.(type) {
	case uint64:
		limit = a.(uint64)
	case uint:
		limit = uint64(a.(uint))
	case int64:
		limit = uint64(a.(int64))
	case string:
		lv, _ := strconv.Atoi(a.(string))
		limit = uint64(lv)
	}
	return
}

func (a *AbstractAdapter) NewColumn(name string, ctype string, options map[string]interface{}) (c Column) {
	c = Column{Name: name, Type: ctype}
	ctype_native_options := a.NativeDatabaseTypes[ctype].(map[string]interface{})

	if as, ok := options["as"]; ok && as.(string) != "" {
		c.As = as.(string)
		return
	} else if as, ok := ctype_native_options["as"]; ok && as.(string) != "" {
		c.As = as.(string)
		return
	}

	var native_limit uint64
	var option_limit uint64

	if flag, ok := ctype_native_options["limit"]; ok {
		native_limit = get_uint64(flag)
		if flag, ok := options["limit"]; ok {
			option_limit = get_uint64(flag)
		}
	}
	if native_limit > 0 {
		if option_limit > 0 {
			c.Limit = option_limit
		} else {
			c.Limit = native_limit
		}
	}

	if flag, ok := options["array"]; ok {
		c.Array = flag.(bool)
	}
	if flag, ok := options["precision"]; ok {
		c.Precision = flag.(int)
	}
	if flag, ok := options["scale"]; ok {
		c.Scale = flag.(int)
	}
	if flag, ok := options["default"]; ok {
		c.Default = flag.(string)
	}

	if flag, ok := options["not_null"]; ok {
		c.NotNull = flag.(bool)
	}
	if flag, ok := options["first"]; ok {
		c.First = flag.(bool)
	}
	if flag, ok := options["after"]; ok {
		c.After = flag.(bool)
	}
	if ctype == "primary_key" {
		c.Primary = true
	}
	if flag, ok := options["primary_key"]; ok {
		c.Primary = flag.(bool)
	}
	return
}

func (c *Column) IsPrimary() bool {
	return c.Primary || c.Type == "primary_key"
}

func (t *Table) PrimaryKeyColumnName() string {
	for _, e := range t.Columns {
		if e.Primary {
			return e.Name
		}
	}
	return ""
}
func (t *Table) Column(name string, ctype string, args ...interface{}) {
	var options map[string]interface{}
	if len(args) > 0 {
		options = args[0].(map[string]interface{})
	}
	if t.PrimaryKeyColumnName() == name {
		log.Fatal("you can't redefine the primary key column '#{name}'. To define a custom primary key, pass { id: false } to create_table.")
		return
	}
	t.Columns = append(t.Columns, t.AbstractAdapter.NewColumn(name, ctype, options))
}

// func (t *Table) PrimaryKey(name string, options map[string]interface{}) {
// 	//options["primary_kay"] = true
// 	t.Column(name, ctype, options)
// }

func (t *Table) PrimaryKey(name string, args ...interface{}) {
	t.Column(name, "primary_key", args...)
}
func (t *Table) String(name string, args ...interface{}) {
	t.Column(name, "string", args...)
}
func (t *Table) Text(name string, args ...interface{}) {
	t.Column(name, "text", args...)
}
func (t *Table) Integer(name string, args ...interface{}) {
	t.Column(name, "integer", args...)
}

func (t *Table) Float(name string, args ...interface{}) {
	t.Column(name, "float", args...)
}
func (t *Table) Decimal(name string, args ...interface{}) {
	t.Column(name, "decimal", args...)
}
func (t *Table) Datetime(name string, args ...interface{}) {
	t.Column(name, "datetime", args...)
}
func (t *Table) Timestamp(name string, args ...interface{}) {
	t.Column(name, "timestamp", args...)
}
func (t *Table) Time(name string, args ...interface{}) {
	t.Column(name, "time", args...)
}
func (t *Table) Date(name string, args ...interface{}) {
	t.Column(name, "date", args...)
}
func (t *Table) Binary(name string, args ...interface{}) {
	t.Column(name, "binary", args...)
}
func (t *Table) Boolean(name string, args ...interface{}) {
	t.Column(name, "boolean", args...)
}
func (t *Table) Timestamps() {
	t.Column("created_at", "timestamp")
	t.Column("updated_at", "timestamp")
}
