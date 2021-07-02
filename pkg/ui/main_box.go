package ui

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/lobshunter86/stop-watch/pkg/core"
	"github.com/lobshunter86/stop-watch/pkg/metrics"
)

func NewUIFromStatuses(statuses map[string]*core.Status, saveStatusHook func()) fyne.Window {
	a := app.New()
	win := a.NewWindow("stop watch")
	win.CenterOnScreen()

	outerBox := container.NewVBox()
	itemList := NewItemList(statuses, outerBox)

	items := make([]*Item, 0, len(statuses))
	for label, status := range statuses {
		ticker := core.NewTicker(time.Second)
		statusIcon := widget.NewIcon(stopImg)
		item := NewItem(label, status, ticker, statusIcon, itemList, saveStatusHook)
		items = append(items, item)
		go item.Start()
	}

	itemList.SetItems(items)
	itemList.listCtner = itemList.ToContainer()

	newItemBotton := NewAddItemBotton(itemList, saveStatusHook)
	outerBox.Add(itemList.listCtner)
	outerBox.Add(newItemBotton)

	win.SetContent(outerBox)
	return win
}

func NewAddItemBotton(itemList *ItemList, saveStatusHook func()) *widget.Button {
	botton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		w := fyne.CurrentApp().NewWindow("enter new item name")
		label := ""
		data := binding.BindString(&label)
		entry := widget.NewEntryWithData(data)
		yesBtn := widget.NewButtonWithIcon("OK", theme.ConfirmIcon(), func() {
			status := core.NewStatus(metrics.Metrics.DurationCount.WithLabelValues(label),
				metrics.Metrics.DurationTotal.WithLabelValues(label), 0)
			statusIcon := widget.NewIcon(stopImg)
			ticker := core.NewTicker(time.Second)
			item := NewItem(label, status, ticker, statusIcon, itemList, saveStatusHook)
			go item.Start()
			err := itemList.AppendItem(item)
			if err != nil {
				fmt.Fprintf(os.Stderr, "AppendItem: %v\n", err)
			}

			w.Close()
		})
		noBtn := widget.NewButtonWithIcon("cancel", theme.CancelIcon(), func() {
			w.Close()
		})
		w.SetContent(container.NewBorder(entry, nil, yesBtn, noBtn))
		w.CenterOnScreen()
		w.Show()
	})

	return botton
}
