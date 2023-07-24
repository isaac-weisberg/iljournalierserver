package main

import (
	"context"
	"fmt"
)

const (
	migrationVersion1 = "1"
)

func migrateDatabase(ctx context.Context, databaseService *databaseService) error {
	transaction, err := databaseService.beginTx(ctx)
	if err != nil {
		return j(err, "transaction creation error")
	}
	defer transaction.rollBack()

	err = transaction.createMigrationsTable()
	if err != nil {
		return j(err, "migrations table creation error")
	}

	v1Migrated, err := transaction.hasVersionBeenMigrated(migrationVersion1)
	if err != nil {
		return j(err, "get hasVersionBeenMigrated failed")
	}

	if !v1Migrated {
		fmt.Println("IlJournalierServer: Migration", migrationVersion1, "migrating...")

		err = transaction.createUsersTable()
		if err != nil {
			return j(err, "create users table error")
		}

		err = transaction.createAccessTokensTable()
		if err != nil {
			return j(err, "create access tokens table failed")
		}
		err = transaction.createMoreMessagesTable()
		if err != nil {
			return j(err, "create more msgs table failed")
		}

		err = transaction.markVersionAsMigrated(migrationVersion1)
		if err != nil {
			return j(err, "markVersionAsMigrated error")
		}
	} else {
		fmt.Println("IlJournalierServer: Migration", migrationVersion1, "already applied")
	}

	err = transaction.commit()

	if err != nil {
		return j(err, "commit failed")
	}

	return nil
}
