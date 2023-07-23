package main

var accessTokensTableReplacePair = ReplacePair{"🪙", "accessTokens"}

func accessTokensSql(sql string) string {
	return replace(sql, accessTokensTableReplacePair)
}

func (transaction transaction) createAccessTokensTable() error {
	sql := replace(
		`CREATE TABLE IF NOT EXISTS 🪙 
		(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL, 
			token TEXT NOT NULL, 
			FOREIGN KEY(userId) REFERENCES 🧑(id)
		)`,
		accessTokensTableReplacePair,
		usersTableReplacePair,
	)

	_, err := transaction.exec(sql)
	if err != nil {
		return j(e("createAccessTokensTable"), err)
	}

	return nil
}
