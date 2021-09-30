package main

import (
	"myhosts/model"
	"myhosts/ui"
	"os"

	"gioui.org/app"
)

func main() {
	model.Init()
	w, aw := ui.CreateWindow()
	shutdown := make(chan int)
	go func() {
		<-shutdown
		model.Close()
		os.Exit(0)
	}()
	go w.Loop(aw, shutdown)
	app.Main()
}
