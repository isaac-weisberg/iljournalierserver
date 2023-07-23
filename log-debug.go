//go:build !prod

package main

import "fmt"

func log(items ...any) {
	fmt.Println(items...)
}
