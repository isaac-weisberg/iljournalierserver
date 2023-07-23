package main

func (transaction transaction) createUsersTable() error {
	errorWrap := createErrorWrapper("tx.createUsersTable")

	sql := "CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, magicKey TEXT NOT NULL UNIQUE)"

	_, err := transaction.exec(sql)
	if err != nil {
		return errorWrap(err)
	}

	return nil
}

func (transaction transaction) createUser(magicKey string) (*int64, error) {
	errorWrap := createErrorWrapper("tx.createUser")

	sql := "INSERT INTO users (magicKey) VALUES (?)"

	result, err := transaction.exec(sql, magicKey)
	if err != nil {
		return nil, errorWrap(err)
	}

	lastIndertedId, err := result.LastInsertId()
	if err != nil {
		return nil, errorWrap(err)
	}

	return &lastIndertedId, nil
}
