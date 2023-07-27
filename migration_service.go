package main

import (
	"context"
	"fmt"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/transaction"
)

const (
	migrationVersion1 = "1"
)

func migrateDatabase(ctx context.Context, databaseService *databaseService) error {
	return beginTxBlockVoid(databaseService, ctx, func(tx *transaction.Transaction) error {
		err := tx.CreateMigrationsTable()
		if err != nil {
			return errors.J(err, "migrations table creation error")
		}

		v1Migrated, err := tx.HasVersionBeenMigrated(migrationVersion1)
		if err != nil {
			return errors.J(err, "get hasVersionBeenMigrated failed")
		}

		if !v1Migrated {
			fmt.Println("IlJournalierServer: Migration", migrationVersion1, "migrating...")

			err = tx.CreateUsersTable()
			if err != nil {
				return errors.J(err, "create users table error")
			}

			err = tx.CreateAccessTokensTable()
			if err != nil {
				return errors.J(err, "create access tokens table failed")
			}

			err = tx.CreateMoreMessagesTable()
			if err != nil {
				return errors.J(err, "create more msgs table failed")
			}

			err = tx.CreateKnownFlagsTable()
			if err != nil {
				return errors.J(err, "createKnownFlagsTable failed")
			}

			err = tx.CreateFlagsTable()
			if err != nil {
				return errors.J(err, "createFlagsTable failed")
			}

			err = tx.MarkVersionAsMigrated(migrationVersion1)
			if err != nil {
				return errors.J(err, "markVersionAsMigrated error")
			}
		} else {
			fmt.Println("IlJournalierServer: Migration", migrationVersion1, "already applied")
		}
		return nil
	})
}
