//go:build !prod

package utils

import "fmt"

func DebugLog(items ...any) {
	fmt.Println(items...)
}
