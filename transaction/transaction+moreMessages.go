package transaction

import "caroline-weisberg.fun/iljournalierserver/errors"

func (transaction *Transaction) CreateMoreMessagesTable() error {
	sql := `
	CREATE TABLE moreMessages (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		userId INTEGER NOT NULL,
		unixSeconds INTEGER NOT NULL,
		message TEXT NON NULL,
		FOREIGN KEY (userId) REFERENCES users(id)
	)`

	_, err := transaction.Exec(sql)
	if err != nil {
		return errors.J(err, "create table failed")
	}
	return nil
}

func (transaction *Transaction) AddMoreMessage(userId int64, unixSeconds int64, msg string) error {
	sql := `
	INSERT INTO moreMessages (userId, unixTime, message) VALUES (?, ?, ?)
	`
	_, err := transaction.Exec(sql, userId, unixSeconds, msg)
	if err != nil {
		return errors.J(err, "insert failed")
	}

	return nil
}
