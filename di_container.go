package main

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/services"
)

type diContainer struct {
	databaseService     *services.DatabaseService
	userService         *userService
	moreMessagesService *moreMessagesService
	flagsService        *flagsService
}

func newDIContainer(ctx context.Context) (*diContainer, error) {
	databaseService, err := services.NewDatabaseService(ctx)

	if err != nil {
		return nil, errors.J(err, "database creation failed")
	}

	err = migrateDatabase(ctx, databaseService)
	if err != nil {
		return nil, errors.J(err, "database migration failed")
	}

	randomIdService := newRandomIdService()

	userService := newUserService(databaseService, &randomIdService)

	moreMessagesService := newMoreMessagesService(databaseService)

	flagsService := newFlagsService(databaseService)

	var di = diContainer{
		databaseService,
		&userService,
		&moreMessagesService,
		&flagsService,
	}
	return &di, nil
}
