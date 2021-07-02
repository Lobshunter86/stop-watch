package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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
	title              *widget.Label
	totalDurationLabel *widget.Label
	startBtn           *widget.Button
	pauseBtn           *widget.Button
	stopBtn            *widget.Button
	deleteBtn          *widget.Button
	statusIcon         *widget.Icon

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

	item := &Item{
		name:    title,
		status:  status,
		stopped: make(chan struct{}, 1),
		ticker:  ticker,

		title:              titleLabel,
		totalDurationLabel: widget.NewLabel(fmt.Sprintf("total: %s", util.FormatDuration(status.Duration))),
		statusIcon:         statusIcon,
		parentItemList:     parentItemList,

		// TODO: stop shall add current duration to total duration & reset current duration
		// pause just simply stops ticker
		startBtn: widget.NewButton("start", func() {
			ticker.Start()
			statusIcon.Resource = startImg
			statusIcon.Refresh()
		}),

		stopBtn: widget.NewButton("stop", func() {
			ticker.Stop()
			statusIcon.Resource = stopImg
			statusIcon.Refresh()
			saveStatusHook()
		}),

		pauseBtn: widget.NewButton("pause", func() {
			// ticker.Stop()
			// statusIcon.Resource = pauseImg
			// statusIcon.Refresh()
		}),
	}

	item.deleteBtn = widget.NewButton("DELETE", func() {
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
	columns := 2
	box := container.NewGridWithRows(
		rows,
		container.NewGridWithColumns(columns, item.title, item.statusIcon),
		item.totalDurationLabel,
		item.startBtn,
		item.pauseBtn,
		item.stopBtn,
		item.deleteBtn,
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
			item.totalDurationLabel.SetText(fmt.Sprintf("total: %s", util.FormatDuration(item.status.Duration)))
		}
	}
}
