package ui

import (
	"errors"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/lobshunter86/stop-watch/pkg/core"
)

var (
	// ErrExist means item already exist
	ErrExist = errors.New("Item already exist")

	// ErrEmptyLabel mean item label is empty
	ErrEmptyLabel = errors.New("Item label is empty")
)

type ItemList struct {
	statuses map[string]*core.Status

	container   *fyne.Container   // parent OuterBox that contains item list
	listCtner   *fyne.Container   // ItemList itself
	scrollCtner *container.Scroll // scrollable list

	// TODO: maybe use map instead of list for items
	// need to keep items in order
	items []*Item // items inside item list
}

func NewItemList(statuses map[string]*core.Status, outerBox *fyne.Container) *ItemList {
	l := &ItemList{
		statuses:  statuses,
		container: outerBox,
	}

	return l
}

func (l *ItemList) SetItems(items []*Item) {
	l.items = items
}

func (l *ItemList) ToContainer() *fyne.Container {
	box := container.NewVBox()
	for _, item := range l.items {
		divider := canvas.NewLine(color.Gray16{0xaff0})
		divider.StrokeWidth = 1
		box.Add(item.itemBox)
		box.Add(divider)
	}

	return box
}

func (l *ItemList) ToScrollContainer() *container.Scroll {
	scroll := container.NewVScroll(l.listCtner)
	return scroll
}

func (l *ItemList) RemoveItem(item *Item) {
	delete(l.statuses, item.name)
	for i, itm := range l.items {
		if item == itm {
			l.items = append(l.items[:i], l.items[i+1:]...)
			break
		}
	}

	l.listCtner.Remove(item.itemBox)
	l.listCtner.Refresh()
	l.container.Refresh()
}

func (l *ItemList) AppendItem(item *Item) error {
	if _, ok := l.statuses[item.name]; ok {
		return ErrExist
	}

	if len(item.name) == 0 {
		return ErrEmptyLabel
	}

	l.statuses[item.name] = item.status
	l.items = append(l.items, item)
	l.listCtner.Add(item.itemBox)
	l.listCtner.Refresh()
	l.container.Refresh()

	return nil
}
