package main

import "github.com/gofrs/uuid"

type UserService struct {
	dbService DatabaseService
}

func NewUserService(dbService DatabaseService) UserService {
	return UserService{dbService: dbService}
}

func (controller *UserService) createUser() (string, error) {
	uniqueId, err := uuid.NewV4()

	if err != nil {
		return "", err
	}

	uuidString := uniqueId.String()

	return uuidString, nil
}
