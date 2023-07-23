package main

func (transaction transaction) createUsersTable() error {
	sql := "CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, magicKey TEXT NOT NULL UNIQUE)"

	_, err := transaction.exec(sql)
	if err != nil {
		return j(e("createUsersTable"), err)
	}

	return nil
}
