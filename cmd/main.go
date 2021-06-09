package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lobshunter86/stop-watch/pkg/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const Argc = 3

var statusFile = "status.json"

func main() {
	if len(os.Args) != Argc {
		fmt.Println(help())
		return
	}

	prometheusPort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(help())
		return
	}

	fmt.Println(version.Version())

	statusFile = os.Args[1]
	status, err := restoreStatus(statusFile)
	if err != nil {
		fmt.Println("restoreStatus: ", err)
		return
	}

	app := app.New()
	w := app.NewWindow("Stopwatch")
	w.SetCloseIntercept(func() {
		err = saveStatus(statusFile, status)
		if err != nil {
			fmt.Println("saveStatus: ", err)
		}
		w.Close()
	})

	label := widget.NewLabel(formatDuration(status.Duration))
	ticker := NewTicker(time.Second)
	startBotton := widget.NewButton("start", ticker.Start)
	stopBotton := widget.NewButton("stop", ticker.Stop)
	go tickLabel(label, &status, ticker)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", prometheusPort), nil) //nolint

	// GUI
	lo := layout.NewBorderLayout(label, nil, startBotton, stopBotton)
	con := container.New(lo, label, startBotton, stopBotton)
	w.SetContent(con)
	w.ShowAndRun()
}

func tickLabel(label *widget.Label, status *Status, ticker *Ticker) {
	for {
		<-ticker.ticker.C
		status.Duration += time.Second
		status.durationCount.WithLabelValues("root").Inc()
		label.SetText(formatDuration(status.Duration))
	}
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	return fmt.Sprintf("%d:%02d:%02d", h, m, d/time.Second)
}

type Status struct {
	durationCount *prometheus.CounterVec
	Duration      time.Duration `json:"duration,omitempty" yaml:"duration"`
}

func restoreStatus(filename string) (Status, error) {
	status := Status{
		durationCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "duration_count",
			Help: "counter of duration seconds",
		},
			[]string{"type"}),
	}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return status, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return Status{}, err
	}

	err = json.Unmarshal(data, &status)
	status.durationCount.WithLabelValues("root").Add(float64(status.Duration / time.Second))
	return status, err
}

func saveStatus(filename string, status Status) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	data, err := json.Marshal(status)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func help() string {
	return ""
}
