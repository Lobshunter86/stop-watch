package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type itemList struct {
	container *fyne.Container
	listCtner *fyne.Container
	items     []*item
	itemBoxs  []*fyne.Container
}

func (l *itemList) toContainer() {
	box := container.NewVBox()
	bigBox := container.NewVBox()
	for _, item := range l.items {
		b := item.ItemBox()
		l.itemBoxs = append(l.itemBoxs, b)

		// delete button
		b.Add(widget.NewButton("DEL", func() {
			w := fyne.CurrentApp().NewWindow("ARE YOU SURE?")
			yesBtn := widget.NewButtonWithIcon("YES", theme.ConfirmIcon(), func() {
				box.Remove(b)
				box.Refresh()
				for i, it := range l.items {
					if item == it {
						l.items = append(l.items[:i], l.items[i+1:]...)
						l.itemBoxs = append(l.itemBoxs[:i], l.itemBoxs[i+1:]...)
					}
				}
				w.Close()
			})
			noBtn := widget.NewButtonWithIcon("no", theme.CancelIcon(), func() {
				w.Close()
			})
			w.SetContent(container.NewBorder(widget.NewLabel("ARE you sure?"), nil, yesBtn, noBtn))
			w.CenterOnScreen()
			w.Show()
		}))
		box.Add(b)
	}
	l.listCtner = box

	// add new item
	newBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		w := fyne.CurrentApp().NewWindow("enter new item name")
		str := ""
		data := binding.BindString(&str)
		entry := widget.NewEntryWithData(data)
		yesBtn := widget.NewButtonWithIcon("OK", theme.ConfirmIcon(), func() {
			// TODO add
			status := widget.NewIcon(stopImg)
			currentTime := widget.NewLabel("current- 3:30:33")
			totalTime := widget.NewLabel("total- 12:22:22")
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
			l.add(str, currentTime, totalTime, startBtn, pauseBtn, stopBtn, status)
			w.Close()
		})
		noBtn := widget.NewButtonWithIcon("cancel", theme.CancelIcon(), func() {
			w.Close()
		})
		w.SetContent(container.NewBorder(entry, nil, yesBtn, noBtn))
		w.CenterOnScreen()
		w.Show()
	})

	bigBox.Add(box)
	bigBox.Add(newBtn)
	l.container = bigBox
}

// add item to itemList's items
func (l *itemList) addItem(itm *item) {
	l.items = append(l.items, itm)
}

// add item to itemList and show it
func (l *itemList) add(title string, currentTime, totalTime *widget.Label, startBtn, pauseBtn, stopBtn *widget.Button,
	status *widget.Icon) {
	itm := NewItem(title, currentTime, totalTime, startBtn, pauseBtn, stopBtn, status)
	l.items = append(l.items, itm)
	b := itm.ItemBox()
	l.itemBoxs = append(l.itemBoxs, b)
	b.Add(widget.NewButton("DEL", func() {
		w := fyne.CurrentApp().NewWindow("ARE YOU SURE?")
		yesBtn := widget.NewButton("YES", func() {
			l.listCtner.Remove(b)
			l.listCtner.Refresh()
			for i, it := range l.items {
				if itm == it {
					l.items = append(l.items[:i], l.items[i+1:]...)
					l.itemBoxs = append(l.itemBoxs[:i], l.itemBoxs[i+1:]...)
				}
			}
			w.Close()
		})
		noBtn := widget.NewButton("no", func() {
			w.Close()
		})
		w.SetContent(container.NewBorder(widget.NewLabel("ARE you sure?"), nil, yesBtn, noBtn))
		w.Show()
	}))
	l.listCtner.Add(b)
	l.listCtner.Refresh()
}
