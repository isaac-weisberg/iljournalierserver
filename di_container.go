package main

import "context"

type diContainer struct {
	databaseService *databaseService
	userService     userService
}

func newDIContainerError(err error) error {
	return j(e("newDIContainer error"), err)
}

func newDIContainer(ctx context.Context) (*diContainer, error) {
	databaseService, err := newDatabaseService(ctx)

	if err != nil {
		return nil, newDIContainerError(err)
	}

	err = migrateDatabase(ctx, *databaseService)
	if err != nil {
		return nil, newDIContainerError(err)
	}

	userService := newUserService(databaseService)

	var di = diContainer{
		databaseService,
		userService,
	}
	return &di, nil
}
