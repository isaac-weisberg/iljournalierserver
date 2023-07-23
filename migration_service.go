package main

import (
	"context"
	"fmt"
)

const (
	migrationVersion1 = "1"
)

var migrateDatabaseError = e("migrateDatabase failed")

func migrateDatabase(ctx context.Context, databaseService databaseService) error {
	transaction, err := databaseService.beginTx(ctx)
	if err != nil {
		return j(migrateDatabaseError, err)
	}
	defer transaction.rollBack()

	err = transaction.createMigrationsTable()
	if err != nil {
		return j(migrateDatabaseError, err)
	}

	v1Migrated, err := transaction.hasVersionBeenMigrated(migrationVersion1)
	if err != nil {
		return j(migrateDatabaseError, err)
	}

	if !v1Migrated {
		fmt.Println("IlJournalierServer: Migration", migrationVersion1, "migrating...")
		err = transaction.markVersionAsMigrated(migrationVersion1)
		if err != nil {
			return j(migrateDatabaseError, err)
		}
	} else {
		fmt.Println("IlJournalierServer: Migration", migrationVersion1, "already applied")
	}

	err = transaction.commit()

	if err != nil {
		return j(migrateDatabaseError, err)
	}

	return nil
}
