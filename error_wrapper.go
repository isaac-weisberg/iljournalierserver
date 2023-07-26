package main

import "errors"

func j(err error, msg string) error {
	return errors.Join(errors.New(msg), err)
}

func js(err error, err2 error) error {
	return errors.Join(err, err2)
}

func e(msg string) error {
	return errors.New(msg)
}

func is(err error, target error) bool {
	return errors.Is(err, target)
}
