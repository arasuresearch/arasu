package abstract

import (
	"strings"
)

func (a *AbstractAdapter) Accept(definition interface{}) (ans string) {
	switch t := definition.(type) {
	case Table:
		ans = a.VisitTable(t)
	case Column:
		ans = a.VisitColumn(t)
	case AlterTable:
		ans = a.VisitAlterTable(t)
	// case AlterColumn:
	// 	ans = a.VisitAlterColumn(t)
	// case DropColumn:
	// 	ans = a.VisitDropColumn(t)

	default:
	}
	return
}

func (a *AbstractAdapter) VisitAddColumn(c Column) string {
	sql := c.TypeToSql(a)
	sql = "ADD " + a.QuoteColumnName(c.Name) + " " + sql
	return a.AddColumnOptions(sql, c)
}

func (a *AbstractAdapter) VisitAlterTable(t AlterTable) string {
	sql := "ALTER TABLE " + a.QuoteTableName(t.Name)
	columns_sql := []string{}
	for _, e := range t.Adds {
		columns_sql = append(columns_sql, a.VisitAddColumn(e))
	}
	return sql + " " + strings.Join(columns_sql, " ")
}

func (a *AbstractAdapter) VisitColumn(c Column) string {
	sql := a.QuoteColumnName(c.Name) + " " + c.TypeToSql(a)
	if !c.Primary {
		sql = a.AddColumnOptions(sql, c)
	}
	return sql
}
func (a *AbstractAdapter) VisitTable(t Table) string {
	sql := "CREATE"
	if t.Temporary {
		sql += " TEMPORARY"
	}
	sql += " TABLE " + a.QuoteTableName(t.Name) + " ("
	//fmt.Println(len(t.Columns), t.Columns)
	if len(t.As) == 0 {
		column_sql := []string{}
		for _, e := range t.Columns {
			column_sql = append(column_sql, a.Accept(e))
		}
		sql += strings.Join(column_sql, ", ")
	} else {
		sql += " AS " + t.As
	}
	sql += ") "
	if flag, ok := t.Options["options"]; ok {
		sql += flag.(string)
	}
	return sql
}

func (a *AbstractAdapter) AddColumnOptions(sql string, c Column) string {
	if len(c.Default) > 0 {
		sql += " DEFAULT " + a.Quote(c.Default)
	}
	if c.NotNull {
		sql += " NOT NULL"
	}
	if a.AutoIncrement() && c.AutoIncrement {
		sql += " AUTO_INCREMENT"
	}
	return sql
}
