package main

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/migrations"
	"caroline-weisberg.fun/iljournalierserver/services"
)

type diContainer struct {
	databaseService     *services.DatabaseService
	userService         *services.UserService
	moreMessagesService *services.MoreMessagesService
	flagsService        *services.FlagsService
}

func newDIContainer(ctx context.Context) (*diContainer, error) {
	databaseService, err := services.NewDatabaseService()

	if err != nil {
		return nil, errors.J(err, "database creation failed")
	}

	err = migrations.MigrateDatabase(ctx, databaseService)
	if err != nil {
		return nil, errors.J(err, "database migration failed")
	}

	randomIdService := services.NewRandomIdService()

	userService := services.NewUserService(databaseService, &randomIdService)

	moreMessagesService := services.NewMoreMessagesService(databaseService)

	flagsService := services.NewFlagsService(databaseService)

	var di = diContainer{
		databaseService,
		&userService,
		&moreMessagesService,
		&flagsService,
	}
	return &di, nil
}
