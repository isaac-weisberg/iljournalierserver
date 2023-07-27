package main

import (
	"context"
	"fmt"
)

const (
	migrationVersion1 = "1"
)

func migrateDatabase(ctx context.Context, databaseService *databaseService) error {
	return beginTxBlockVoid(databaseService, ctx, func(tx *transaction) error {
		err := tx.createMigrationsTable()
		if err != nil {
			return j(err, "migrations table creation error")
		}

		v1Migrated, err := tx.hasVersionBeenMigrated(migrationVersion1)
		if err != nil {
			return j(err, "get hasVersionBeenMigrated failed")
		}

		if !v1Migrated {
			fmt.Println("IlJournalierServer: Migration", migrationVersion1, "migrating...")

			err = tx.createUsersTable()
			if err != nil {
				return j(err, "create users table error")
			}

			err = tx.createAccessTokensTable()
			if err != nil {
				return j(err, "create access tokens table failed")
			}

			err = tx.createMoreMessagesTable()
			if err != nil {
				return j(err, "create more msgs table failed")
			}

			err = tx.createKnownFlagsTable()
			if err != nil {
				return j(err, "createKnownFlagsTable failed")
			}

			err = tx.createFlagsTable()
			if err != nil {
				return j(err, "createFlagsTable failed")
			}

			err = tx.markVersionAsMigrated(migrationVersion1)
			if err != nil {
				return j(err, "markVersionAsMigrated error")
			}
		} else {
			fmt.Println("IlJournalierServer: Migration", migrationVersion1, "already applied")
		}
		return nil
	})
}
