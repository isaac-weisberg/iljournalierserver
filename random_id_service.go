package main

import "github.com/gofrs/uuid"

type randomIdService struct {
}

func newRandomIdService() randomIdService {
	return randomIdService{}
}

func randomIdServiceError(err error) error {
	return j(e("randomIdServiceError"), err)
}

func (randomIdService randomIdService) generateRandomId() (*string, error) {
	uniqueId, err := uuid.NewV4()

	if err != nil {
		return nil, randomIdServiceError(err)
	}

	str := uniqueId.String()

	return &str, nil
}
