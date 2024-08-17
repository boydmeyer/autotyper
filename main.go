package main

import (
	"os"
	"time"

	g "xabbo.b7c.io/goearth"
	"xabbo.b7c.io/goearth/shockwave/in"
)

var ext = g.NewExt(g.ExtInfo{
	Title:       "Autotyper",
	Description: "The AutoTyper app features an intuitive, easy-to-use interface that simplifies repetitive messaging by automatically sending predefined messages at customizable intervals.",
	Version:     "2.0.0",
	Author:      "Nanobyte",
})

var quitChan = make(chan struct{})
var at *AutoTyper

func main() {
	at = NewAutoTyper(ext)

	ext.Intercept(in.TRADE_COMPLETED, in.TRADE_COMPLETED_2).With(OnTradeCompleted)

	go func() {
		ext.RunE()
	}()

	go func() {
		<-quitChan
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	at.Run()
}

func OnTradeCompleted(e *g.Intercept) {
	at.onTradeCompleted()
}
