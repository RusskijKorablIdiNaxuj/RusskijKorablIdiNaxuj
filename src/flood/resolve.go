//go:build !wasm
// +build !wasm

package flood

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	skipHostResolution = "skip"
)

const requestTimeout = time.Millisecond * 300

func additionalHeaders(request *http.Request) {
}

func (t *Target) isHostResolved() bool {
	t.Lock()
	defer t.Unlock()

	if time.Since(t.lastResolved) > time.Minute*5 {
		t.resolvedAddress = t.resolve(t.address)
		t.lastResolved = time.Now()
	}

	return t.resolvedAddress != nil
}

func (t *Target) replaceWithResolvedIP(addr string) string {
	// return addr
	t.RLock()
	defer t.RUnlock()

	if t.resolvedAddress == nil {
		return ""
	}
	if t.resolvedAddress[0] == skipHostResolution {
		return addr
	}

	url, err := url.Parse(addr)
	if err != nil {
		t.exitCh <- struct{}{}
		return ""
	}
	url.Host = t.resolvedAddress[rand.Intn(len(t.resolvedAddress))]

	return url.String()
}

func (t *Target) resolve(addr string) []string {
	url, err := url.Parse(addr)
	if err != nil {
		t.exitCh <- struct{}{}
		return nil
	}

	if net.ParseIP(url.Hostname()) != nil {
		return []string{skipHostResolution}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addrs, err := net.DefaultResolver.LookupHost(ctx, url.Hostname())
	if err != nil {
		return nil
	}

	validIP := []string{}
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip == nil {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			continue
		}
		if ip.IsLoopback() || ip.IsPrivate() {
			continue
		}
		validIP = append(validIP, ip.String())
	}

	return validIP
}
