package services

import (
	"caroline-weisberg.fun/iljournalierserver/errors"
	"github.com/gofrs/uuid"
)

type RandomIdService struct {
}

func NewRandomIdService() RandomIdService {
	return RandomIdService{}
}

func (randomIdService *RandomIdService) GenerateRandomId() (*string, error) {
	uniqueId, err := uuid.NewV4()

	if err != nil {
		return nil, errors.J(err, "uuid creation error")
	}

	str := uniqueId.String()

	return &str, nil
}
