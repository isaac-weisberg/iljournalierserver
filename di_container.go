package main

import "context"

type diContainer struct {
	databaseService *databaseService
	userService     userService
}

var newDIContainerError = e("newDIContainer error")

func newDIContainer(ctx context.Context) (*diContainer, error) {
	databaseService, err := newDatabaseService(ctx)

	if err != nil {
		return nil, j(newDIContainerError, err)
	}

	err = migrateDatabase(ctx, *databaseService)
	if err != nil {
		return nil, j(newDIContainerError, err)
	}

	userService := newUserService(databaseService)

	var di = diContainer{
		databaseService,
		userService,
	}
	return &di, nil
}
