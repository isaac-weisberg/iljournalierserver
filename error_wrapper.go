package main

import "errors"

func j(err error, msg string) error {
	return errors.Join(errors.New(msg), err)
}
