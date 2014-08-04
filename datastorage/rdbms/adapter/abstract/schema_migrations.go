package abstract

func (a *AbstractAdapter) CreateSchemaMigration() error {
	txn, err := a.DB.Begin()
	if err != nil {
		return err
	}
	a.Tx = txn
	err = a.CreateTable("schema_migrations", map[string]interface{}{"id": false}, func(t *Table) {
		t.String("version", map[string]interface{}{"limit": 255, "not_null": true})
	})
	if err != nil {
		txn.Rollback()
		return err
	}
	txn.Commit()
	return nil
}

// func (a *AbstractAdapter) DropSchemaMigration() error {
// 	txn, err := a.DB.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	a.Tx = txn
// 	err = a.DropTable("schema_migrations")
// 	if err != nil {
// 		txn.Rollback()
// 		return err
// 	}
// 	txn.Commit()
// 	return nil
// }
func (a *AbstractAdapter) InsertIntoSchemaMigration(version string) error {
	sql := "INSERT INTO " + a.Quote("schema_migrations") + " (" + a.Quote("version") + ") VALUES (" + a.SingleQuote(version) + ")"
	if _, err := a.ExecWithTxn(sql); err != nil {
		return err
	}
	return nil
}
func (a *AbstractAdapter) DeleteFromSchemaMigration(version string) error {
	sql := "DELETE FROM " + a.Quote("schema_migrations") + " WHERE " + a.Quote("schema_migrations") + "." + a.Quote("version") + " = " + a.SingleQuote(version)
	if _, err := a.ExecWithTxn(sql); err != nil {
		return err
	}
	return nil
}

func (a *AbstractAdapter) GetAllSchemaMigration() ([]string, error) {
	ok, err := a.TableExists("schema_migrations")
	if err != nil {
		return nil, err
	}
	if !ok {
		if err := a.CreateSchemaMigration(); err != nil {
			return nil, err
		}
	}

	sql := "SELECT " + a.Quote("schema_migrations") + ".* FROM " + a.Quote("schema_migrations")
	return a.QueryRowsFirstColumnStringArray(sql)
}
