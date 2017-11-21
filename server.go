package main

import (
	"fmt"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"time"
)

func HelloWorld(histogram *prometheus.HistogramVec) http.HandlerFunc {
	return func (res http.ResponseWriter, req *http.Request) {
		start := time.Now()


		defer func() {
			duration := time.Since(start)
			histogram.WithLabelValues(fmt.Sprintf("%d", http.StatusOK)).Observe(duration.Seconds())
		}()

		log.Println("Hello World!")
		fmt.Fprint(res, "Hello World")
	}

}


func main() {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "hello_world_seconds",
		Help: "Time taken to return hello world",
	}, []string{"code"})

	prometheus.Register(histogram)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", HelloWorld(histogram))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
