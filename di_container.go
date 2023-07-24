package main

import "context"

type diContainer struct {
	databaseService *databaseService
	userService     userService
}

func newDIContainer(ctx context.Context) (*diContainer, error) {
	databaseService, err := newDatabaseService(ctx)

	if err != nil {
		return nil, j(err, "database creation failed")
	}

	err = migrateDatabase(ctx, *databaseService)
	if err != nil {
		return nil, j(err, "database migration failed")
	}

	randomIdService := newRandomIdService()

	userService := newUserService(databaseService, &randomIdService)

	var di = diContainer{
		databaseService,
		userService,
	}
	return &di, nil
}
