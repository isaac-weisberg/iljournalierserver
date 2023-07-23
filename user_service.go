package main

import "github.com/gofrs/uuid"

type userService struct {
	dbService *databaseService
}

func newUserService(dbService *databaseService) userService {
	return userService{dbService: dbService}
}

func (userService *userService) createUser() (string, error) {
	uniqueId, err := uuid.NewV4()

	if err != nil {
		return "", err
	}

	uuidString := uniqueId.String()

	return uuidString, nil
}
