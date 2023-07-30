package transaction

import (
	"strings"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/models"
)

func (transaction *Transaction) CreateMoreMessagesTable() error {
	query := `
	CREATE TABLE moreMessages (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		userId INTEGER NOT NULL,
		unixSeconds INTEGER NOT NULL,
		message TEXT NON NULL,
		FOREIGN KEY (userId) REFERENCES users(id)
	)`

	_, err := transaction.Exec(query)
	if err != nil {
		return errors.J(err, "create table failed")
	}
	return nil
}

func (transaction *Transaction) AddMoreMessages(userId int64, addMessageRequests []models.AddMessageRequest) error {
	if len(addMessageRequests) == 0 {
		return errors.E("no messages to add")
	}

	var queryBuilder strings.Builder
	var args = make([]any, 0, len(addMessageRequests)*3)

	var firstAddMessageRequest = addMessageRequests[0]
	queryBuilder.WriteString("INSERT INTO moreMessages (userId, unixSeconds, message) VALUES (?, ?, ?)")
	args = append(args, userId, firstAddMessageRequest.UnixSeconds, firstAddMessageRequest.Message)

	var remainingMessageRequests = addMessageRequests[1:]
	for _, messageRequest := range remainingMessageRequests {
		queryBuilder.WriteString(", (?, ?, ?)")
		args = append(args, userId, messageRequest.UnixSeconds, messageRequest.Message)
	}

	var sql = queryBuilder.String()

	_, err := transaction.Exec(sql, args...)
	if err != nil {
		return errors.J(err, "insert failed")
	}

	return nil
}
