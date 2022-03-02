//go:build wasm
// +build wasm

// There is no possibility to resolve hostnames from webasm...
package flood

import (
	"net/http"
	"time"
)

const requestTimeout = time.Millisecond * 50

func additionalHeaders(request *http.Request) {
	request.Header.Add("js.fetch:mode", "no-cors")
	request.Header.Add("Access-Control-Allow-Origin", "*")
}

func (t *Target) isHostResolved() bool {
	return true
}

func (t *Target) replaceWithResolvedIP(addr string) string {
	return addr
}
