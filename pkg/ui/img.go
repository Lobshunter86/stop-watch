package ui

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

var (

	//go:embed img/start.png
	startImgContent []byte
	//go:embed img/pause.png
	pauseImgContent []byte
	//go:embed img/stop.png
	stopImgContent []byte
	//go:embed img/stopBtn.png
	stopBtnImgContent []byte

	startImg   = fyne.NewStaticResource("start.png", startImgContent)
	pauseImg   = fyne.NewStaticResource("pause.png", pauseImgContent)
	stopImg    = fyne.NewStaticResource("stop.png", stopImgContent)
	stopBtnImg = fyne.NewStaticResource("stopBtn.png", stopBtnImgContent)
)
