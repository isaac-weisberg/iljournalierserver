package main

import "context"

type diContainer struct {
	databaseService     *databaseService
	userService         *userService
	moreMessagesService *moreMessagesService
}

func newDIContainer(ctx context.Context) (*diContainer, error) {
	databaseService, err := newDatabaseService(ctx)

	if err != nil {
		return nil, j(err, "database creation failed")
	}

	err = migrateDatabase(ctx, databaseService)
	if err != nil {
		return nil, j(err, "database migration failed")
	}

	randomIdService := newRandomIdService()

	userService := newUserService(databaseService, &randomIdService)

	moreMessagesService := newMoreMessagesService(databaseService)

	var di = diContainer{
		databaseService,
		&userService,
		&moreMessagesService,
	}
	return &di, nil
}
