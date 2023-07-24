package main

import "github.com/gofrs/uuid"

type randomIdService struct {
}

func newRandomIdService() randomIdService {
	return randomIdService{}
}

func (randomIdService *randomIdService) generateRandomId() (*string, error) {
	uniqueId, err := uuid.NewV4()

	if err != nil {
		return nil, j(err, "uuid creation error")
	}

	str := uniqueId.String()

	return &str, nil
}
