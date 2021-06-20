package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"

	"github.com/lobshunter86/stop-watch/pkg/core"
	"github.com/lobshunter86/stop-watch/pkg/metrics"
	"github.com/lobshunter86/stop-watch/pkg/ui"
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
			fmt.Fprintln(os.Stderr, "NewFileStore: ", err)
			return err
		}

		durations := statusStore.GetAll()
		statuses := make(map[string]*core.Status, len(durations))
		for label, duration := range durations {
			statuses[label] = core.NewStatus(metrics.Metrics.DurationCount.WithLabelValues(label),
				metrics.Metrics.DurationTotal.WithLabelValues(label), duration)

			statuses[label].TotalCounter.Add(float64(duration / time.Second))
		}

		http.Handle("/metrics", promhttp.Handler())
		go http.ListenAndServe(fmt.Sprintf(":%d", prometheusPort), nil) //nolint

		w := ui.NewUIFromStatuses(statuses)
		onClose := func() {
			s := make(map[string]time.Duration)
			for label, status := range statuses {
				s[label] = status.TotalDuration
			}

			err = statusStore.Save(s)
			if err != nil {
				fmt.Println("saveStatus: ", err)
			}
			w.Close()
		}

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
		w.SetCloseIntercept(onClose)
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
