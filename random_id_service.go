package main

import "github.com/gofrs/uuid"

type randomIdService struct {
}

func newRandomIdService() randomIdService {
	return randomIdService{}
}

func (randomIdService randomIdService) generateRandomId() (*string, error) {
	wrapError := createErrorWrapper("randomIdServiceError")

	uniqueId, err := uuid.NewV4()

	if err != nil {
		return nil, wrapError(err)
	}

	str := uniqueId.String()

	return &str, nil
}
