package abstract

import (
	"fmt"
	"github.com/arasuresearch/arasu/lib"
	"strings"
)

func (a *AbstractAdapter) Tables(name, database, like string) ([]string, error) {
	sql := "SHOW TABLES"
	if len(database) > 0 {
		sql += " IN " + a.SingleQuote(database)
	}
	if len(like) > 0 {
		sql += " LIKE " + a.SingleQuote(like)
	}
	names, err := a.QueryRowsFirstColumnStringArray(sql)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func (a *AbstractAdapter) TableExists(name string) (bool, error) {
	if len(name) == 0 {
		return false, fmt.Errorf("table name should't be empty")
	}
	tables, err := a.Tables("", "", name)
	if err != nil {
		return false, err
	}
	if lib.StringArrayContains(tables, name) {
		return true, nil
	}

	if strings.Contains(name, ".") {
		names := strings.Split(name, ".")
		l := len(names)
		name = names[l-1]
		database := names[l-2]

		tables, err := a.Tables("", database, name)
		if err != nil {
			return false, err
		}
		if lib.StringArrayContains(tables, name) {
			return true, nil
		}
	}
	return false, nil
}
