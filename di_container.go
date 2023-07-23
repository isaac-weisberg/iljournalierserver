package main

import "context"

type diContainer struct {
	databaseService *databaseService
	userService     userService
}

func newDIContainer(ctx context.Context) (*diContainer, error) {
	databaseService, err := newDatabaseService(ctx)

	if err != nil {
		return nil, err
	}

	userService := newUserService(databaseService)

	var di = diContainer{
		databaseService,
		userService,
	}
	return &di, nil
}
