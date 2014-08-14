package abstract

import (
	"encoding/json"
	"github.com/arasuresearch/arasu/lib/stringer"
	"log"
	"reflect"
	"strings"
)

//CREATE  INDEX `index_people_on_name_and_age`  ON `people` (`name`, `age`)

func (a *AbstractAdapter) IndexExists(name string) bool {
	return stringer.Contains([]string{}, name)
}
func (a *AbstractAdapter) CreateIndex(table_name string, column_names []string, options map[string]interface{}) error {

	index_name, index_type, index_columns, index_options, _, _ := a.AddIndexOptions(table_name, column_names, options)
	sql := "CREATE " + index_type + " INDEX " + a.QuoteColumnName(index_name) + " ON " +
		a.QuoteTableName(table_name) + " (" + index_columns + ") "
	if len(index_options) > 0 {
		io, _ := json.Marshal(index_options)
		sql += string(io)
	}
	if _, err := a.Tx.Exec(sql); err != nil {
		return err
	}
	return nil
}
func (a *AbstractAdapter) DropIndex(table_name string, column_names []string, options map[string]interface{}) error {
	index_name, _, _, _, _, _ := a.AddIndexOptions(table_name, column_names, options)
	sql := "DROP INDEX " + a.QuoteColumnName(index_name) + " ON " + a.QuoteTableName(table_name)
	if _, err := a.Tx.Exec(sql); err != nil {
		return err
	}
	return nil
}

// func (a *AbstractAdapter) DropIndex(table_name string, column_names []string, options map[string]interface{}) error {
// 	index_name := a.index_name_for_remove(table_name, options)
// 	sql := "DROP INDEX " + a.QuoteColumnName(index_name) + " ON " + a.QuoteTableName(table_name)
// 	if _, err := a.Tx.Exec(sql); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (a *AbstractAdapter) AlterIndex(table_name string, old_name string, new_name string) {
// 	var old_index_def *Index
// 	for _, e := range a.Indexes(table_name) {
// 		if e.Name == old_name {
// 			old_index_def = &e
// 		}
// 	}
// 	if old_index_def == nil {
// 		return
// 	}
// 	a.CreateIndex(table_name, old_index_def.ColumnNames, map[string]interface{}{"name": new_name, "unique": old_index_def.Unique})
// 	a.DropIndex(table_name, map[string]string{"name": old_name})
// }

// func (a *AbstractAdapter) RenameIndex(table_name string, old_name string, new_name string) {
// 	var old_index_def *Index
// 	for _, e := range a.Indexes(table_name) {
// 		if e.Name == old_name {
// 			old_index_def = &e
// 		}
// 	}
// 	if old_index_def == nil {
// 		return
// 	}
// 	a.CreateIndex(table_name, old_index_def.ColumnNames, map[string]interface{}{"name": new_name, "unique": old_index_def.Unique})
// 	a.DropIndex(table_name, map[string]string{"name": old_name})
// }

func (a *AbstractAdapter) Indexes(table_name string) []Index {
	return []Index{}
}

func (a *AbstractAdapter) IndexNameExists(table_name string, index_name string, not_implemented bool) bool {
	if m := reflect.ValueOf(a).MethodByName("Indexes"); m.Interface() == nil {
		return not_implemented
	} else {
		for _, e := range a.Indexes(table_name) {
			if e.Name == index_name {
				return true
			}
		}
	}
	return false
}

func (a *AbstractAdapter) IndexName(table_name string, options interface{}) (index_name string) {
	switch options.(type) {
	case map[string]interface{}:
		mapped_options := options.(map[string]interface{})
		if column_names, exists := mapped_options["column"]; exists {
			index_name = "index_" + table_name + "_on_" + strings.Join(column_names.([]string), "_and_")
		} else if name, exists := mapped_options["name"]; exists {
			index_name = name.(string)
		} else {
			log.Fatal("You must specify the index name")
		}
	default:
		index_name = a.IndexName(table_name, map[string]interface{}{"column": options.([]string)})
	}
	return
}

func (a *AbstractAdapter) AddIndexOptions(table_name string, column_names []string, options map[string]interface{}) (index_name string, index_type string, index_columns string, index_options string, algorithm string, using string) {
	if name, ok := options["name"]; ok {
		index_name = name.(string)
	} else if len(table_name) > 0 && len(column_names) > 0 {
		index_name = "index_" + table_name + "_on_" + strings.Join(column_names, "_and_")
	} else {
		log.Fatal("You must specify the index name")
	}
	if flag, ok := options["unique"]; ok {
		if flag.(bool) == true {
			index_type = "UNIQUE"
		}
	}
	if flag, ok := options["type"]; ok {
		index_type = flag.(string)
	}
	if flag, ok := options["using"]; ok {
		using = flag.(string)
	}
	if a.SupportPartialIndex() {
		if flag, ok := options["where"]; ok {
			index_options = " WHERE " + flag.(string)
		}
	}
	if len(index_name) > a.MaxIndexLength {
		log.Fatal("Index name '#{index_name}' on table '#{table_name}' is too long; the limit is #{max_index_length} characters")
	}
	if a.IndexNameExists(table_name, index_name, false) {
		log.Fatal("Index name '#{index_name}' on table '#{table_name}' already exists")
	}
	index_columns = strings.Join(a.QuotedColumnsForIndex(column_names, options), ", ")

	return
}

func (a *AbstractAdapter) AddIndexSortOrder(option_strings map[string]string, column_names []string, options map[string]interface{}) map[string]string {
	if order, ok := options["order"]; ok {
		for _, e := range column_names {
			option_strings[e] += " " + strings.ToUpper(order.(string))
		}
	}
	return option_strings
}

func (a *AbstractAdapter) QuotedColumnsForIndex(column_names []string, options map[string]interface{}) []string {
	option_strings := map[string]string{}
	for _, e := range column_names {
		option_strings[e] = " "
	}
	if a.SupportsIndexSortOrder() {
		option_strings = a.AddIndexSortOrder(option_strings, column_names, options)
	}
	for i, e := range column_names {
		column_names[i] = a.QuoteColumnName(e) + option_strings[e]
	}
	return column_names
}

func (a *AbstractAdapter) index_name_for_remove(table_name string, options map[string]interface{}) (index_name string) {
	index_name = a.IndexName(table_name, options)
	if !a.IndexNameExists(table_name, index_name, true) {
		if _, exists := options["name"]; exists {
			options_without_column := options
			delete(options_without_column, "column")
			index_name_without_column := a.IndexName(table_name, options_without_column)
			if a.IndexNameExists(table_name, index_name_without_column, false) {
				return index_name_without_column
			}
		}
		log.Fatal("Index name '" + index_name + "' on table '" + table_name + "' does not exist")
	}
	return
}
