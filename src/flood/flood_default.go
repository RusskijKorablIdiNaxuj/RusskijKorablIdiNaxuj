//go:build test
// +build test

//-go:build !http && !sun
//- +build !http,!sun

package flood

import (
	"errors"
	"math/rand"
	"time"
)

func (t *Target) perform(addr string) error {
	time.Sleep(time.Millisecond + time.Duration(rand.Int63n(1000)))
	if rand.Intn(10) <= 5 {
		return errors.New("test")
	}

	return nil
}
