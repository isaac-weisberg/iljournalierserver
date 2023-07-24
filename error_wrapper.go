package main

import "errors"

func j(err error, msg string) error {
	return errors.Join(errors.New(msg), err)
}

func e(msg string) error {
	return errors.New(msg)
}

func is(err error, target error) bool {
	return errors.Is(err, target)
}
