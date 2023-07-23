package main

func (transaction transaction) createAccessTokensTable() error {
	sql := `CREATE TABLE IF NOT EXISTS accessTokens 
		(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL, 
			token TEXT NOT NULL UNIQUE, 
			FOREIGN KEY(userId) REFERENCES users(id)
		)`

	_, err := transaction.exec(sql)
	if err != nil {
		return j(e("createAccessTokensTable"), err)
	}

	return nil
}
