package flood

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

type Target struct {
	address   string
	port      int
	randomize bool

	requestsPerSecond float64
	errorsPerSecond   float64
	requests          int64
	errors            int64

	requestCh chan string

	urls []string

	client *http.Client
}

func New(addr string) Target {
	tr := &http.Transport{
		MaxIdleConns:       10000,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	return Target{
		address:   strings.Trim(addr, " \r\n\t"),
		port:      80,
		randomize: true,
		client:    client,
	}
}

func (t *Target) Name() string {
	return t.address
}

func (t *Target) Run(ctx context.Context, N, maxRPS int, progress func(requests, errors int64)) {
	t.requestCh = make(chan string, N)
	timer := time.NewTicker(time.Second)
	timerGen := time.NewTicker(time.Second / time.Duration(maxRPS))
	defer close(t.requestCh)
	defer timer.Stop()
	defer timerGen.Stop()

	for i := 0; i < N; i++ {
		go t.flood(ctx)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-timerGen.C:
			t.requestCh <- t.generate()
		case <-timer.C:
			requests := float64(atomic.SwapInt64(&t.requests, 0))
			errors := float64(atomic.SwapInt64(&t.errors, 0))
			t.requestsPerSecond = (t.requestsPerSecond*3 + requests) / 4
			t.errorsPerSecond = (t.errorsPerSecond*3 + errors) / 4
			progress(int64(t.requestsPerSecond), int64(t.errorsPerSecond))
		}
	}
}

func (t *Target) flood(ctx context.Context) {
	for addr := range t.requestCh {
		atomic.AddInt64(&t.requests, 1)
		if t.perform(ctx, addr) != nil {
			atomic.AddInt64(&t.errors, 1)
		}
	}
}

func (t *Target) generate() string {
	if strings.HasPrefix(t.address, "https://") || strings.HasPrefix(t.address, "http://") {
		return t.address
	}

	protoArr := []string{"https", "http"}
	port := t.port
	proto := protoArr[rand.Intn(len(protoArr))]

	url := ""

	return fmt.Sprintf("%s://%s:%d/%s", proto, t.address, port, url)
}
