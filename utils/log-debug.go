//go:build !prod

package utils

import "fmt"

func Log(items ...any) {
	fmt.Println(items...)
}
