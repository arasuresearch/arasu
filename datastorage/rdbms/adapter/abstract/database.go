package abstract

func (a *AbstractAdapter) CreateDatabase() error {
	if err := a.InitWithoutDb(); err != nil {
		return err
	}
	charset := DEFAULT_CHAR_SET
	if flag, ok := a.Options["charset"]; ok {
		charset = flag[0]
	}
	sql := "CREATE DATABASE " + a.Quote(a.DbName) + " DEFAULT CHARACTER SET " + a.Quote(charset)
	if flag, ok := a.Options["collation"]; ok {
		sql += " COLLATE " + a.Quote(flag[0])
	}
	if _, err := a.ExecWithTxn(sql); err != nil {
		return err
	}
	return nil
}
func (a *AbstractAdapter) DropDatabase() error {
	if err := a.InitWithoutDb(); err != nil {
		return err
	}
	//query := "DROP DATABASE IF EXISTS " + a.Quote(a.DbName)
	query := "DROP DATABASE " + a.Quote(a.DbName)

	if _, err := a.ExecWithTxn(query); err != nil {
		return err
	}
	return nil
}
func (a *AbstractAdapter) CurrentDatabase() (string, error) {
	res, err := a.QueryRowsFirstColumnStringArray("SELECT DATABASE() as db")
	if err != nil {
		return "", err
	}
	return res[0], nil
}
