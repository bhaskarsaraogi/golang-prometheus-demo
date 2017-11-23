package main

import (
	"fmt"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"strconv"
	"runtime"
)

const INSTANCE = "prod0"

func HelloWorld(histogram *prometheus.HistogramVec, counter *prometheus.CounterVec, gauge *prometheus.GaugeVec, summary *prometheus.SummaryVec) http.HandlerFunc {
	return func (res http.ResponseWriter, req *http.Request) {
		timer := prometheus.NewTimer(histogram.With(prometheus.Labels{
			"status": strconv.Itoa(http.StatusOK),
			"instance": INSTANCE}))

		summaryTimer := prometheus.NewTimer(summary.With(prometheus.Labels{
			"status": strconv.Itoa(http.StatusOK),
			"instance": INSTANCE}))

		defer func() {

			var m runtime.MemStats

			runtime.ReadMemStats(&m)
			gauge.With(prometheus.Labels{"instance":INSTANCE}).Set(float64(m.Alloc / 1024))

			counter.With(prometheus.Labels{"instance":INSTANCE}).Inc()

			timer.ObserveDuration()
			summaryTimer.ObserveDuration()
		}()

		log.Println("Hello World!")
		fmt.Fprint(res, "Hello World")
	}

}


func main() {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "hello_world_response_time",
		Help: "Time taken to return hello world",
	}, []string{"status", "instance"})

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hello_world_requests_count",
		Help: "Number of hellow world requests recvd",
	}, []string{"instance"})

	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "used_memory",
		Help: "Head used memory footprint in kilobytes",
	}, []string{"instance"})

	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "yet_another_response_time_tracker",
		Help: "Summary of request response time",
	}, []string{"status", "instance"})

	prometheus.Register(histogram)
	prometheus.Register(counter)
	prometheus.Register(gauge)
	prometheus.Register(summary)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", HelloWorld(histogram, counter, gauge, summary))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
