package main

import (
	"fmt"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"time"
	"strconv"
)

const INSTANCE = "prod0"

func HelloWorld(histogram *prometheus.HistogramVec, counter *prometheus.CounterVec) http.HandlerFunc {
	return func (res http.ResponseWriter, req *http.Request) {
		start := time.Now()


		defer func() {
			duration := time.Since(start)
			histogram.With(prometheus.Labels{
				"status": strconv.Itoa(http.StatusOK),
				"instance": INSTANCE}).Observe(duration.Seconds())
			counter.With(prometheus.Labels{"instance":INSTANCE}).Inc()
		}()

		log.Println("Hello World!")
		fmt.Fprint(res, "Hello World")
	}

}


func main() {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "hello_world_seconds",
		Help: "Time taken to return hello world",
	}, []string{"status", "instance"})

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hello_world_requests_count",
		Help: "Number of hellow world requests recvd",
	}, []string{"instance"})

	prometheus.Register(histogram)
	prometheus.Register(counter)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", HelloWorld(histogram, counter))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
