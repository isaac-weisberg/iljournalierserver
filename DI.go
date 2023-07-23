package main

type DIContainer struct {
	databaseService DatabaseService
	userService     UserService
}

func NewDIContainer() (DIContainer, error) {
	databaseService, err := NewDatabaseService()

	if err != nil {
		return DIContainer{}, err
	}

	userService := NewUserService(databaseService)

	var diContainer = DIContainer{
		databaseService,
		userService,
	}
	return diContainer, nil
}
