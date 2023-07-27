package main

import (
	"caroline-weisberg.fun/iljournalierserver/errors"
	"github.com/gofrs/uuid"
)

type randomIdService struct {
}

func newRandomIdService() randomIdService {
	return randomIdService{}
}

func (randomIdService *randomIdService) generateRandomId() (*string, error) {
	uniqueId, err := uuid.NewV4()

	if err != nil {
		return nil, errors.J(err, "uuid creation error")
	}

	str := uniqueId.String()

	return &str, nil
}
