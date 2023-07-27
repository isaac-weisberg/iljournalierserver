package main

import (
	"strings"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

func (transaction *transaction) createFlagsTable() error {
	query := `CREATE TABLE flags (
		id INTEGER NOT NULL PRIMARY KEY,
		unixSeconds INTEGER NOT NULL,
		flagId INTEGER NOT NULL,
		FOREIGN KEY (flagId) REFERENCES knownFlags(id)
	)`

	_, err := transaction.exec(query)
	if err != nil {
		return errors.J(err, "create table failed")
	}

	return nil
}

func (transaction *transaction) markFlags(requests []markFlagRequest) error {
	if len(requests) == 0 {
		return nil
	}

	firstRequest := requests[0]
	remainingRequests := requests[1:]

	builder := strings.Builder{}
	args := []any{firstRequest.UnixSeconds, firstRequest.FlagId}

	builder.WriteString("INSERT INTO flags (unixSeconds, flagId) VALUES (?, ?)")

	for _, request := range remainingRequests {
		builder.WriteString(", (?, ?)")
		args = append(args, request.UnixSeconds, request.FlagId)
	}

	resultingQuery := builder.String()

	_, err := transaction.exec(resultingQuery, args)
	if err != nil {
		return errors.J(err, "inserting failed")
	}

	return nil
}
