package main

func createErrorWrapper(message string) func(err error) error {
	return func(err error) error {
		return j(e(message), err)
	}
}
