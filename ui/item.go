package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type item struct {
	title       *widget.Label
	currentTime *widget.Label
	totalTime   *widget.Label
	startBtn    *widget.Button
	pauseBtn    *widget.Button
	stopBtn     *widget.Button
	status      *widget.Icon
}

func (itm *item) ItemBox() *fyne.Container {

	row := 3
	box := container.NewGridWithRows(
		row,
		container.NewGridWithColumns(2, itm.title, itm.status),
		itm.currentTime,
		itm.totalTime,
		itm.startBtn,
		itm.pauseBtn,
		layout.NewSpacer(),
		itm.stopBtn,
	)

	return box
}

func NewItem(title string, currentTime, totalTime *widget.Label, startBtn, pauseBtn, stopBtn *widget.Button, status *widget.Icon) *item {
	// make title center
	var align fyne.TextAlign = 1
	style := &fyne.TextStyle{
		Bold:      true,
		Italic:    true,
		Monospace: false,
	}
	titleLabel := widget.NewLabelWithStyle(title, align, *style)
	itm := &item{
		title:       titleLabel,
		currentTime: currentTime,
		totalTime:   totalTime,
		startBtn:    startBtn,
		pauseBtn:    pauseBtn,
		stopBtn:     stopBtn,
		status:      status,
	}
	return itm
}
