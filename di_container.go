package main

import "context"

type diContainer struct {
	databaseService *databaseService
	userService     userService
}

func newDIContainer(ctx context.Context) (*diContainer, error) {
	wrapError := createErrorWrapper("newDIContainer error")
	databaseService, err := newDatabaseService(ctx)

	if err != nil {
		return nil, wrapError(err)
	}

	err = migrateDatabase(ctx, *databaseService)
	if err != nil {
		return nil, wrapError(err)
	}

	randomIdService := newRandomIdService()

	userService := newUserService(databaseService, &randomIdService)

	var di = diContainer{
		databaseService,
		userService,
	}
	return &di, nil
}
