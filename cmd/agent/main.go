package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type gauge float64
type counter int64

type Metric struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge

	PollCount   counter
	RandomValue gauge
}

func NewMonitor() {
	var m Metric
	var rtm runtime.MemStats
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second
	pollTimer := time.NewTicker(pollInterval)
	reportTimer := time.NewTicker(reportInterval)

	time.AfterFunc(0, func() {
		for {
			<-reportTimer.C
			fmt.Println("AfterFunc", m)
		}
	})

	for {
		<-pollTimer.C

		runtime.ReadMemStats(&rtm)

		m.Alloc = gauge(rtm.Alloc)
		m.BuckHashSys = gauge(rtm.BuckHashSys)
		m.Frees = gauge(rtm.Frees)
		m.GCCPUFraction = gauge(rtm.GCCPUFraction)
		m.GCSys = gauge(rtm.GCSys)
		m.HeapAlloc = gauge(rtm.HeapAlloc)
		m.HeapIdle = gauge(rtm.HeapIdle)
		m.HeapInuse = gauge(rtm.HeapInuse)
		m.HeapObjects = gauge(rtm.HeapObjects)
		m.HeapReleased = gauge(rtm.HeapReleased)
		m.HeapSys = gauge(rtm.HeapSys)
		m.LastGC = gauge(rtm.LastGC)
		m.Lookups = gauge(rtm.Lookups)
		m.MCacheInuse = gauge(rtm.MCacheInuse)
		m.MCacheSys = gauge(rtm.MCacheSys)
		m.MSpanInuse = gauge(rtm.MSpanInuse)
		m.MSpanSys = gauge(rtm.MSpanSys)
		m.Mallocs = gauge(rtm.Mallocs)
		m.NextGC = gauge(rtm.NextGC)
		m.NumForcedGC = gauge(rtm.NumForcedGC)
		m.NumGC = gauge(rtm.NumGC)
		m.OtherSys = gauge(rtm.OtherSys)
		m.PauseTotalNs = gauge(rtm.PauseTotalNs)
		m.StackInuse = gauge(rtm.StackInuse)
		m.StackSys = gauge(rtm.StackSys)
		m.Sys = gauge(rtm.Sys)
		m.TotalAlloc = gauge(rtm.TotalAlloc)

		m.PollCount += 1
		m.RandomValue = gauge(rand.Intn(1000000000000000))

		// // Just encode to json and print
		b, _ := json.Marshal(m)
		fmt.Println(string(b))
	}
}

func main() {
	NewMonitor()
}
