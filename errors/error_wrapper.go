package errors

import "errors"

func J(err error, msg string) error {
	return errors.Join(errors.New(msg), err)
}

func Js(err error, err2 error) error {
	return errors.Join(err, err2)
}

func E(msg string) error {
	return errors.New(msg)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}
