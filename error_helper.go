package main

import "errors"

func e(string string) error {
	return errors.New(string)
}

func j(err1 error, err2 error) error {
	return errors.Join(err1, err2)
}
