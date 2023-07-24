package main

func (transaction *transaction) createMoreMessagesTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS moreMessages (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		userId INTEGER NOT NULL, 
		message TEXT NON NULL,
		FOREIGN KEY (userId) REFERENCES users(id)
	)`

	_, err := transaction.exec(sql)
	if err != nil {
		return j(err, "create table failed")
	}
	return nil
}

func (transaction *transaction) addMoreMessage(userId int64, msg string) error {
	sql := `
	INSERT INTO moreMessages (userId, message) VALUES (?, ?)
	`
	_, err := transaction.exec(sql, userId, msg)
	if err != nil {
		return j(err, "insert failed")
	}

	return nil
}
