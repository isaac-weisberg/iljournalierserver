package main

var usersTableReplacePair = ReplacePair{"🧑", "users"}

func usersTableSql(sql string) string {
	return replace(sql, usersTableReplacePair)
}

func (transaction transaction) createUsersTable() error {
	sql := usersTableSql("CREATE TABLE IF NOT EXISTS 🧑 (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, hash TEXT NOT NULL, salt TEXT NOT NULL)")

	_, err := transaction.exec(sql)
	if err != nil {
		return j(e("createUsersTable"), err)
	}

	return nil
}
