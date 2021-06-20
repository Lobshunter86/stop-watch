package ui

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

var (
	//go:embed img/start.png
	startImgContent []byte
	//go:embed img/stop.png
	stopImgContent []byte

	startImg = fyne.NewStaticResource("start.png", startImgContent)
	stopImg  = fyne.NewStaticResource("stop.png", stopImgContent)
)
