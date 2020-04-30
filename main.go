package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const RemoteLabel = "remote"

type Size struct {
	Count int `json:"count"`
	Bytes int `json:"bytes"`
}

func recordMetrics() {
	go func() {
		for {
			var rcloneSizeCmd = exec.Command("rclone", "size", remote, "--json")

			var out bytes.Buffer
			var stderr bytes.Buffer

			rcloneSizeCmd.Stdout = &out
			rcloneSizeCmd.Stderr = &stderr

			err := rcloneSizeCmd.Run()

			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
				connectionErrorsGauge.With(prometheus.Labels{RemoteLabel: remote}).Inc()
			} else {
				var response Size
				json.Unmarshal(out.Bytes(), &response)
				connectionErrorsGauge.With(prometheus.Labels{RemoteLabel: remote}).Set(float64(0))
				countGauge.With(prometheus.Labels{RemoteLabel: remote}).Set(float64(response.Count))
				sizeGauge.With(prometheus.Labels{RemoteLabel: remote}).Set(float64(response.Bytes))
			}
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}()
}

var (
	countGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rclone_size_exporter_objects_number",
		Help: "Number of elements on the remote volume.",
	},
		[]string{RemoteLabel})

	sizeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rclone_size_exporter_total_size",
		Help: "Size of the remote volume.",
	},
		[]string{RemoteLabel})

	connectionErrorsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rclone_size_exporter_connection_error",
		Help: "Flags if there are current issues connecting to rclone.",
	},
		[]string{RemoteLabel})

	remote = os.Getenv("REMOTE")

	delay = getEnvAsInt("DELAY", 300)
)

func main() {

	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
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
