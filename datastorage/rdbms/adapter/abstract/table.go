package abstract

func (a *AbstractAdapter) Columns(table_name string) []string {
	// for _, e := range a.TableNames {
	// 	if e==table_name {

	// 	}
	// }
	// return lib.StringArrayContains(a.IndexNames, tablename)
	return []string{}
}

func (a *AbstractAdapter) NewTable(name string, options map[string]interface{}) Table {
	table := Table{Name: name, AbstractAdapter: a, Options: options}
	if flag, ok := options["temporary"]; ok {
		table.Temporary = flag.(bool)
	}
	if flag, ok := options["as"]; ok {
		table.As = flag.(string)
	}
	return table
}

func (a *AbstractAdapter) CreateTable(name string, args ...interface{}) error {
	//options map[string]interface{}, callback func(t *Table)
	var callback func(t *Table)
	var options = map[string]interface{}{}
	var primary_key string = "id"
	var option_id bool = true

	if len(args) == 1 {
		callback = args[0].(func(t *Table))
	} else if len(args) > 1 {
		options = args[0].(map[string]interface{})
		callback = args[1].(func(t *Table))
	}
	if _, ok := options["options"]; !ok {
		options["options"] = "ENGINE=InnoDB"
	}

	table := a.NewTable(name, options)
	if pk, ok := options["primary_key"]; ok {
		primary_key = pk.(string)
	}
	if id, ok := options["id"]; ok {
		option_id = id.(bool)
	}

	if len(table.As) == 0 {
		if option_id {
			table.PrimaryKey(primary_key, options)
		}
		callback(&table)
	}
	if flag, ok := options["force"]; ok {
		if value, ok := flag.(bool); ok && value {
			ok, err := a.TableExists(table.Name)
			if err != nil {
				return err
			} else if ok {
				if err := a.DropTable(table.Name, options); err != nil {
					return err
				}
			}

		}
	}
	if _, err := a.Tx.Exec(a.Accept(table)); err != nil {
		return err
	}
	for _, e := range table.Indexes {
		if err := a.CreateIndex(table.Name, e.ColumnNames, options); err != nil {
			//a.DropTable(name)
			return err
		}
	}
	return nil
}
func (a *AbstractAdapter) DropTable(table_name string, args ...interface{}) error {
	_, err := a.Tx.Exec("DROP TABLE " + a.QuoteTableName(table_name))
	return err
}

func (a *AbstractAdapter) AlterTable(name string, args ...interface{}) error {

	var callback func(t *AlterTable)
	var options = map[string]interface{}{}

	if len(args) == 1 {
		callback = args[0].(func(t *AlterTable))
	} else if len(args) > 1 {
		options = args[0].(map[string]interface{})
		callback = args[1].(func(t *AlterTable))
	}
	_, _ = options, callback
	return nil

}

// func (a *AbstractAdapter) BulkChangeTable(table_name string, operations) {
//         sqls = operations.mapk do |command, args|
//           table, arguments = args.shift, args
//           method = :"#{command}_sql"
//           if respond_to?(method, true)
//             send(method, table, *arguments)
//           else
//             raise "Unknown method called : #{method}(#{arguments.inspect})"
//           end
//         end.flatten.join(", ")
//         execute("ALTER TABLE #{quote_table_name(table_name)} #{sqls}")
// }

// func (a *AbstractAdapter) ChangeTable(table_name string, options map[string]string) {
//         if a.SupportsBulkAlter() && options["bulk"]=="true" {
//           yield update_table_definition(table_name, recorder)
//           bulk_change_table(table_name, recorder.commands)
//         } else {
//           yield update_table_definition(table_name, self)

//         }
//       }
