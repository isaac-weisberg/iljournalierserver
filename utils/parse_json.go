package utils

import (
	"encoding/json"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"github.com/go-playground/validator/v10"
)

func ParseJson[R interface{}](input []byte) (*R, error) {
	var parsedBody R
	err := json.Unmarshal(input, &parsedBody)
	if err != nil {
		return nil, errors.J(err, "parsing json failed")
	}

	var validate = validator.New()

	err = validate.Struct(parsedBody)
	if err != nil {
		return nil, errors.J(err, "validating struct failed")
	}

	return &parsedBody, nil
}
