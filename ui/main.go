package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

var (
	stopImg, _  = fyne.LoadResourceFromPath("img/stop.png")
	startImg, _ = fyne.LoadResourceFromPath("img/start.png")
	pauseImg, _ = fyne.LoadResourceFromPath("img/pause.png")
)

func main() {
	a := app.New()
	win := a.NewWindow("title")

	title := "paint"
	title2 := "work"
	currentTime := widget.NewLabel("current- 3:30:33")
	totalTime := widget.NewLabel("total- 12:22:22")

	status := widget.NewIcon(stopImg)

	startBtn := widget.NewButton("start", func() {
		status.Resource = startImg
		status.Refresh()
	})
	pauseBtn := widget.NewButton("pause", func() {
		status.Resource = pauseImg
		status.Refresh()
	})
	stopBtn := widget.NewButton("stop", func() {
		status.Resource = stopImg
		status.Refresh()
	})
	itm := NewItem(title, currentTime, totalTime, startBtn, pauseBtn, stopBtn, status)
	itm2 := NewItem(title2, currentTime, totalTime, startBtn, pauseBtn, stopBtn, status)
	itmList := itemList{}
	itmList.addItem(itm)
	itmList.addItem(itm2)
	itmList.toContainer()
	win.SetContent(itmList.container)
	win.CenterOnScreen()
	win.ShowAndRun()
}
