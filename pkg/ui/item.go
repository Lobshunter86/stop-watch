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
	title                *widget.Label
	currentDurationLabel *widget.Label
	totalDurationLabel   *widget.Label
	startBtn             *widget.Button
	pauseBtn             *widget.Button
	stopBtn              *widget.Button
	deleteBtn            *widget.Button
	statusIcon           *widget.Icon

	status *core.Status
	ticker *core.Ticker
	name   string
}

func NewItem(title string, status *core.Status, ticker *core.Ticker, statusIcon *widget.Icon, parentItemList *ItemList) *Item {
	var align fyne.TextAlign = 1
	style := &fyne.TextStyle{
		Bold:      true,
		Italic:    true,
		Monospace: false,
	}
	titleLabel := widget.NewLabelWithStyle(title, align, *style)

	item := &Item{
		name:   title,
		status: status,
		ticker: ticker,

		title:                titleLabel,
		currentDurationLabel: widget.NewLabel(fmt.Sprintf("current: %s", util.FormatDuration(0))),
		totalDurationLabel:   widget.NewLabel(fmt.Sprintf("total: %s", util.FormatDuration(status.Duration))),
		statusIcon:           statusIcon,
		parentItemList:       parentItemList,

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
	rows := 3
	columns := 2
	box := container.NewGridWithRows(
		rows,
		container.NewGridWithColumns(columns, item.title, item.statusIcon),
		item.currentDurationLabel,
		item.totalDurationLabel,
		item.startBtn,
		item.pauseBtn,
		layout.NewSpacer(),
		item.stopBtn,
		item.deleteBtn,
	)

	return box
}

func (item *Item) Start() {
	// FIXME
	// should exit when item is deleted
	// otherwise there is reource leak

	for {
		<-item.ticker.C
		item.status.Duration += time.Second
		item.status.Counter.Inc()
		item.status.TotalCounter.Inc()
		item.totalDurationLabel.SetText(fmt.Sprintf("total: %s", util.FormatDuration(item.status.Duration)))
	}
}
