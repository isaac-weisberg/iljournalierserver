package main

import "errors"

func createErrorWrapper(message string) func(err error) error {
	return func(err error) error {
		return errors.Join(errors.New(message), err)
	}
}
