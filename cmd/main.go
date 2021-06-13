package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/lobshunter86/stop-watch/pkg/core"
	"github.com/lobshunter86/stop-watch/pkg/util"
	"github.com/lobshunter86/stop-watch/pkg/version"
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
	statusStore, err := core.NewFileStore(statusFile)
	if err != nil {
		fmt.Println("NewFileStore: ", err)
		return
	}
	statusTime := statusStore.GetAll()["root"]
	status := NewStatus(durationCount.WithLabelValues("root"), statusTime)
	status.Counter.Add(float64(statusTime / time.Second))

	app := app.New()
	w := app.NewWindow("Stopwatch")

	onClose := func() {
		err = statusStore.Save(map[string]time.Duration{"root": status.Duration})
		if err != nil {
			fmt.Println("saveStatus: ", err)
		}
		w.Close()
	}
	w.SetCloseIntercept(onClose)

	label := widget.NewLabel(util.FormatDuration(status.Duration))
	ticker := core.NewTicker(time.Second)
	startBotton := widget.NewButton("start", ticker.Start)
	stopBotton := widget.NewButton("stop", ticker.Stop)
	go tickLabel(label, &status, ticker)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", prometheusPort), nil) //nolint

	// catch signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		sig := <-sc
		fmt.Printf("Got signal: [%s], shutting down\n", sig.String())
		onClose()
	}()

	// GUI
	lo := layout.NewBorderLayout(label, nil, startBotton, stopBotton)
	con := container.New(lo, label, startBotton, stopBotton)
	w.SetContent(con)
	w.ShowAndRun()
}

func tickLabel(label *widget.Label, status *Status, ticker *core.Ticker) {
	for {
		<-ticker.C
		status.Duration += time.Second
		status.Counter.Inc()
		label.SetText(util.FormatDuration(status.Duration))
	}
}

type Status struct {
	Counter  prometheus.Counter
	Duration time.Duration `json:"duration,omitempty" yaml:"duration"`
}

func NewStatus(c prometheus.Counter, duration time.Duration) Status {
	return Status{
		Counter:  c,
		Duration: duration,
	}
}

func help() string {
	return ""
}
