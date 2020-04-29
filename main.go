package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Size struct {
	Count int `json:"count"`
	Bytes int `json:"bytes"`
}

func recordMetrics() {
	go func() {
		for {
			out, err := rcloneSizeCmd.CombinedOutput()

			if err != nil {
				connectionErrorsGauge.Inc()
			} else {
				var response Size
				json.Unmarshal(out, &response)
				connectionErrorsGauge.Set(float64(0))
				countGauge.Set(float64(response.Count))
				sizeGauge.Set(float64(response.Bytes))
			}
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}()
}

var (
	countGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "rclone_objects_number",
		Help: "Number of elements on the remote volume.",
	})

	sizeGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "rclone_total_size",
		Help: "Size of the remote volume.",
	})

	connectionErrorsGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "rclone_connection_error",
		Help: "Flags if there are current issues connecting to rclone.",
	})

	rcloneSizeCmd = exec.Command("rclone", "--config=/config/rclone.conf", "size", "jim:", "--json")
	delay         = getEnvAsInt("DELAY", 30)
)

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}

//Utilities

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
