package main

import (
	"context"
	"fmt"
)

const (
	migrationVersion1 = "1"
)

func migrateDatabaseError(err error) error {
	return j(e("migrateDatabase failed"), err)
}

func migrateDatabase(ctx context.Context, databaseService databaseService) error {

	transaction, err := databaseService.beginTx(ctx)
	if err != nil {
		return migrateDatabaseError(err)
	}
	defer transaction.rollBack()

	err = transaction.createMigrationsTable()
	if err != nil {
		return migrateDatabaseError(err)
	}

	v1Migrated, err := transaction.hasVersionBeenMigrated(migrationVersion1)
	if err != nil {
		return migrateDatabaseError(err)
	}

	if !v1Migrated {
		fmt.Println("IlJournalierServer: Migration", migrationVersion1, "migrating...")

		err = transaction.createUsersTable()
		if err != nil {
			return migrateDatabaseError(err)
		}

		err = transaction.createAccessTokensTable()
		if err != nil {
			return migrateDatabaseError(err)
		}

		err = transaction.markVersionAsMigrated(migrationVersion1)
		if err != nil {
			return migrateDatabaseError(err)
		}
	} else {
		fmt.Println("IlJournalierServer: Migration", migrationVersion1, "already applied")
	}

	err = transaction.commit()

	if err != nil {
		return migrateDatabaseError(err)
	}

	return nil
}
