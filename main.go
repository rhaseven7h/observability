package main

import (
	_ "expvar"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
)

const (
	BindAddress = "0.0.0.0"
	BindPort    = 6060
	LogPath     = "/var/log/app.kibanatest.log"
)

func main() {
	logFile, err := os.OpenFile(LogPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logFile.Close()
	}()
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.InfoLevel)
	logrusLogger.SetFormatter(&ecslogrus.Formatter{
		DataKey: "labels",
	})
	logrusLogger.SetOutput(logFile)
	logrusLogger.ReportCaller = true
	lg := logrus.NewEntry(logrusLogger).
		WithField("app", "kibanatest").
		WithField("owner", "gmm")

	bindAddress := fmt.Sprintf("%s:%d", BindAddress, BindPort)
	lg.
		WithField("bind_address", bindAddress).
		Infoln("STARTING KIBANA TEST")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(lg *logrus.Entry) {
		for {
			lg.
				WithField("action", "legal").
				WithField("to_legal", rand.Uint32()).
				Info("legal update")
			if rand.Intn(100) > 80 {
				lg.
					WithField("action", "excess").
					WithField("excess_level", rand.Uint32()).
					Error("excess error")
			}
			time.Sleep(time.Second)
		}
	}(lg)

	eventsProcessed := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "kibanatestns",
		Subsystem: "kibanatestss",
		Name:      "kibana_test_event_processed_total",
		Help:      "The total number of events_processed",
		ConstLabels: map[string]string{
			"ktlabel1": "value1",
			"ktlabel2": "value2",
			"ktlabel3": "value3",
		},
	})
	go func(eventsProcessed prometheus.Counter) {
		for {
			eventsProcessed.Add(float64(1 + rand.Intn(10)))
			time.Sleep(200 * time.Millisecond)
		}
	}(eventsProcessed)

	http.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lg.
			WithField("action", "metrics").
			WithField("version", "v0.0.1").
			Info("metrics delivery")
		promhttp.Handler().ServeHTTP(w, r)
	}))
	err = http.ListenAndServe(bindAddress, nil)
	if err != nil {
		panic(err)
	}
}
