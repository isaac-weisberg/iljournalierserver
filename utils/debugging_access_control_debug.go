//go:build !prod

package utils

import "net/http"

var debuggingOrigins = [...]string{
	"https://localhost:9000",
	"http://localhost:9000",
}

func WriteDebugCorsHeader(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	var originToWriteIn string
	if origin == "" {
		originToWriteIn = debuggingOrigins[0]
	} else {
		for _, debuggingOrigin := range debuggingOrigins {
			if origin == debuggingOrigin {
				originToWriteIn = debuggingOrigin
				goto write // I'm sorry lol
			}
		}
		originToWriteIn = origin
	}

write:
	w.Header().Add("Access-Control-Allow-Origin", originToWriteIn)
}
