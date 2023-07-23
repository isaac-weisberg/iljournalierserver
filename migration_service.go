package main

import (
	"context"
	"database/sql"
)

var createMigrationsTable = "CREATE TABLE IF NOT EXISTS migrations (version INT NOT NULL, PRIMARY KEY (version))"

func migrateDatabase(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, createMigrationsTable)
	if err != nil {
		return err
	}

	return nil
}
