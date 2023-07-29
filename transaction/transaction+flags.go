package transaction

import (
	"strings"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

func (transaction *Transaction) CreateFlagsTable() error {
	query := `CREATE TABLE flags (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		unixSeconds INTEGER NOT NULL,
		flagId INTEGER NOT NULL,
		FOREIGN KEY (flagId) REFERENCES knownFlags(id)
	)`

	_, err := transaction.Exec(query)
	if err != nil {
		return errors.J(err, "create table failed")
	}

	return nil
}

type MarkFlagRequest struct {
	UnixSeconds int64
	FlagId      int64
}

func (transaction *Transaction) MarkFlags(requests []MarkFlagRequest) error {
	if len(requests) == 0 {
		return nil
	}

	firstRequest := requests[0]
	remainingRequests := requests[1:]

	builder := strings.Builder{}

	var args = make([]any, 0, len(requests))
	args = append(args, firstRequest.UnixSeconds, firstRequest.FlagId)

	builder.WriteString("INSERT INTO flags (unixSeconds, flagId) VALUES (?, ?)")

	for _, request := range remainingRequests {
		builder.WriteString(", (?, ?)")
		args = append(args, request.UnixSeconds, request.FlagId)
	}

	resultingQuery := builder.String()

	_, err := transaction.Exec(resultingQuery, args...)
	if err != nil {
		return errors.J(err, "inserting failed")
	}

	return nil
}
