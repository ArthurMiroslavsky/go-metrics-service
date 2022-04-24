package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type gauge float64
type counter int64

type mapOfRequests map[string]*http.Request
type metricRequests func() mapOfRequests
type gaugesMap map[string]gauge
type counterMap map[string]counter

const pollInterval = 2 * time.Second
const reportInterval = 10 * time.Second
const clientTimeout = time.Second * 5
const url = "http://127.0.0.1:8080/"

var gauges map[string]gauge
var counters counterMap

func sendRequest(client *http.Client, fn metricRequests) {
	metricRequestsMap := fn()

	for _, value := range metricRequestsMap {
		response, err := client.Do(value)

		if err != nil {
			log.Fatal("Error when sending the request", err)
		}

		defer response.Body.Close()
	}

}

func gaugeMetricRequests(gauges gaugesMap) mapOfRequests {
	m := make(mapOfRequests)

	for key, value := range gauges {
		reqUrl := fmt.Sprintf("%s/update/gauge/%s/%f", url, key, value)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBufferString(""))

		if err != nil {
			log.Printf("An error occurred while sending metrics (metricType: gauge, metricName: %s, metricValue: %f", key, value)
		}

		request.Header.Set("Content-Type", "text/plain")

		m[key] = request
	}

	return m

}

func counterMetricRequests(counters counterMap) mapOfRequests {
	m := make(mapOfRequests)

	for key, value := range counters {
		reqUrl := fmt.Sprintf("%s/update/counter/%s/%d", url, key, value)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBufferString(""))

		if err != nil {
			log.Printf("An error occurred while sending metrics (metricType: counter, metricName: %s, metricValue: %d", key, value)
		}

		request.Header.Set("Content-Type", "text/plain")

		m[key] = request
	}

	return m
}

func NewMonitor(ctx context.Context) {
	gauges := make(map[string]gauge)
	counters := make(counterMap)

	var rtm runtime.MemStats
	var counter counter

	pollTimer := time.NewTicker(pollInterval)
	reportTimer := time.NewTicker(reportInterval)

	transport := &http.Transport{
		MaxIdleConns: 30,
	}

	client := &http.Client{
		Timeout:   clientTimeout,
		Transport: transport,
	}

	for {
		select {
		case <-ctx.Done():
			pollTimer.Stop()
			reportTimer.Stop()
			log.Println("shutting down server gracefully")
			return
		default:
			time.AfterFunc(0, func() {
				for {
					<-reportTimer.C
					countersReqFn := func() mapOfRequests { return counterMetricRequests(counters) }
					gaugesReqFn := func() mapOfRequests { return gaugeMetricRequests(gauges) }

					sendRequest(client, countersReqFn)
					sendRequest(client, gaugesReqFn)
				}
			})

			for {
				<-pollTimer.C

				counter++

				runtime.ReadMemStats(&rtm)

				gauges["Alloc"] = gauge(rtm.Alloc)
				gauges["BuckHashSys"] = gauge(rtm.BuckHashSys)
				gauges["Frees"] = gauge(rtm.Frees)
				gauges["GCCPUFraction"] = gauge(rtm.GCCPUFraction)
				gauges["GCSys"] = gauge(rtm.GCSys)
				gauges["HeapAlloc"] = gauge(rtm.HeapAlloc)
				gauges["HeapIdle"] = gauge(rtm.HeapIdle)
				gauges["HeapInuse"] = gauge(rtm.HeapInuse)
				gauges["HeapObjects"] = gauge(rtm.HeapObjects)
				gauges["HeapReleased"] = gauge(rtm.HeapReleased)
				gauges["HeapSys"] = gauge(rtm.HeapSys)
				gauges["LastGC"] = gauge(rtm.LastGC)
				gauges["Lookups"] = gauge(rtm.Lookups)
				gauges["MCacheInuse"] = gauge(rtm.MCacheInuse)
				gauges["MCacheSys"] = gauge(rtm.MCacheSys)
				gauges["MSpanInuse"] = gauge(rtm.MSpanInuse)
				gauges["MSpanSys"] = gauge(rtm.MSpanSys)
				gauges["Mallocs"] = gauge(rtm.Mallocs)
				gauges["NextGC"] = gauge(rtm.NextGC)
				gauges["NumForcedGC"] = gauge(rtm.NumForcedGC)
				gauges["NumGC"] = gauge(rtm.NumGC)
				gauges["OtherSys"] = gauge(rtm.OtherSys)
				gauges["PauseTotalNs"] = gauge(rtm.PauseTotalNs)
				gauges["StackInuse"] = gauge(rtm.StackInuse)
				gauges["StackSys"] = gauge(rtm.StackSys)
				gauges["Sys"] = gauge(rtm.Sys)
				gauges["TotalAlloc"] = gauge(rtm.TotalAlloc)

				gauges["RandomValue"] = gauge(rand.Intn(1000000000000000))

				counters["PollCount"] = counter
			}
		}
	}

}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	NewMonitor(ctx)
}
