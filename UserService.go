package main

import "github.com/gofrs/uuid"

type UserService struct {
}

func (controller *UserService) createUser() (string, error) {
	uniqueId, err := uuid.NewV4()

	if err != nil {
		return "", err
	}

	uuidString := uniqueId.String()

	return uuidString, nil
}
