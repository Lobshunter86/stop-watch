package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/lobshunter86/stop-watch/pkg/core"
	"github.com/lobshunter86/stop-watch/pkg/util"
)

type Item struct {
	// UI
	parentItemList *ItemList
	itemBox        *fyne.Container // represents item itself

	// components in itemBox
	title *widget.Label
	// test currnentDuration
	isPause              *bool
	isStart              *bool
	currentDuration      *time.Duration
	currentDurationLabel *widget.Label
	totalDurationLabel   *widget.Label
	startBtn             *widget.Button
	pauseBtn             *widget.Button
	stopBtn              *widget.Button
	deleteBtn            *widget.Button
	statusIcon           *widget.Icon

	status  *core.Status
	ticker  *core.Ticker
	stopped chan struct{}
	name    string
}

func NewItem(title string, status *core.Status, ticker *core.Ticker, statusIcon *widget.Icon, parentItemList *ItemList, saveStatusHook func()) *Item {
	var align fyne.TextAlign = 1
	style := &fyne.TextStyle{
		Bold:      true,
		Italic:    true,
		Monospace: false,
	}
	titleLabel := widget.NewLabelWithStyle(title, align, *style)

	// test currnentDuration
	isPause := false
	isStart := false
	var currentDuration time.Duration = 0
	currentDurationLabel := widget.NewLabel(fmt.Sprintf("C: %s", util.FormatDuration(currentDuration)))

	item := &Item{
		name:    title,
		status:  status,
		stopped: make(chan struct{}, 1),
		ticker:  ticker,

		title:              titleLabel,
		totalDurationLabel: widget.NewLabel(fmt.Sprintf("T: %s", util.FormatDuration(status.Duration))),
		// test currnentDuration
		isPause:              &isPause,
		isStart:              &isStart,
		currentDuration:      &currentDuration,
		currentDurationLabel: currentDurationLabel,
		statusIcon:           statusIcon,
		parentItemList:       parentItemList,

		// TODO: stop shall add current duration to total duration & reset current duration
		// pause just simply stops ticker
		startBtn: widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
			// test currnentDuration
			ticker.Start()
			statusIcon.Resource = startImg
			statusIcon.Refresh()
			if !isStart {
				if isPause {
					isPause = false
				} else {
					currentDuration = 0
				}
				isStart = true
				currentDurationLabel.SetText(fmt.Sprintf("C: %s", util.FormatDuration(currentDuration)))
			}
		}),

		stopBtn: widget.NewButtonWithIcon("", stopBtnImg, func() {
			ticker.Stop()
			isStart = false
			isPause = false
			statusIcon.Resource = stopImg
			statusIcon.Refresh()
			saveStatusHook()
		}),

		pauseBtn: widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {
			ticker.Stop()
			isPause = true
			statusIcon.Resource = pauseImg
			statusIcon.Refresh()
			saveStatusHook()
		}),
	}

	item.deleteBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		w := fyne.CurrentApp().NewWindow("ARE YOU SURE?")
		yesBtn := widget.NewButtonWithIcon("YES", theme.ConfirmIcon(), func() {
			item.parentItemList.RemoveItem(item)
			item.Stop()
			w.Close()
		})
		noBtn := widget.NewButtonWithIcon("no", theme.CancelIcon(), func() {
			w.Close()
		})
		w.SetContent(container.NewBorder(widget.NewLabel("ARE you sure?"), nil, yesBtn, noBtn))
		w.CenterOnScreen()
		w.Show()
	})

	item.itemBox = item.toContainer()
	return item
}

func (item *Item) toContainer() *fyne.Container {
	rows := 2
	// columns := 2
	// box := container.NewGridWithRows(
	// 	rows,
	// 	// container.NewGridWithColumns(columns, item.title, item.statusIcon),
	// 	item.title,
	// 	item.statusIcon,
	// 	currentLabel,
	// 	totalLabel,
	// 	item.currentDurationLabel,
	// 	item.totalDurationLabel,
	// 	item.startBtn,
	// 	item.stopBtn,
	// 	item.deleteBtn,
	// )
	box := container.NewVBox(
		container.NewGridWithRows(rows, item.title, item.statusIcon,
			item.currentDurationLabel, item.totalDurationLabel),
		container.NewHBox(
			layout.NewSpacer(),
			item.startBtn,
			item.pauseBtn,
			item.stopBtn,
			item.deleteBtn,
			layout.NewSpacer(),
		),
	)

	return box
}

func (item *Item) Stop() {
	item.stopped <- struct{}{}
}

func (item *Item) Start() {
	for {
		select {
		case <-item.stopped:
			return

		case <-item.ticker.C:
			item.status.Duration += time.Second
			item.status.Counter.Inc()
			item.status.TotalCounter.Inc()
			item.totalDurationLabel.SetText(fmt.Sprintf("T: %s", util.FormatDuration(item.status.Duration)))
			// test currnentDuration
			*item.currentDuration += time.Second
			item.currentDurationLabel.SetText(fmt.Sprintf("C: %s", util.FormatDuration(*item.currentDuration)))
		}
	}
}
