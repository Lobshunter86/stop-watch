package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"

	"github.com/lobshunter86/stop-watch/pkg/core"
	"github.com/lobshunter86/stop-watch/pkg/util"
	"github.com/lobshunter86/stop-watch/pkg/version"
)

var statusFile string
var prometheusPort int

func init() {
	rootCmd.Flags().StringVarP(&statusFile, "statusfile", "s", "status.json", "path to status file")
	rootCmd.Flags().IntVarP(&prometheusPort, "port", "p", 9111, "prometheus metrics port")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "stop-watch",
	Short:   "A Stop Watch",
	Version: version.Version(),
	RunE: func(cmd *cobra.Command, args []string) error {
		statusStore, err := core.NewFileStore(statusFile)
		if err != nil {
			fmt.Println("NewFileStore: ", err)
			return err
		}
		statusTime := statusStore.GetAll()["root"]
		status := core.NewStatus(metrics.durationCount.WithLabelValues("root"), metrics.durationTotal.WithLabelValues("root"), statusTime)
		status.TotalCounter.Add(float64(statusTime / time.Second))

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
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func tickLabel(label *widget.Label, status *core.Status, ticker *core.Ticker) {
	for {
		<-ticker.C
		status.Duration += time.Second
		status.Counter.Inc()
		status.TotalCounter.Inc()
		label.SetText(util.FormatDuration(status.Duration))
	}
}
