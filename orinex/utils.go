package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type Utils struct{}

type CpuInfo struct {
	User   float64
	System float64
	Alarm  bool
}

type Statistic struct {
	TotalConnections   int32 `json:"total connections"`
	RequestRate        int64 `json:"request rate"`
	CurrentConnections int32 `json:"remain jobs"`
	FailConnections    int32 `json:"failed connections"`
	ProcessJob         int32 `json:"processed jobs"`
	ProcessRate        int64 `json:"process rate"`
	SystemAlarm        bool  `json:"alarm"`
}

func (c CpuInfo) String() string {
	return fmt.Sprintf("(User: %.2f) (System: %.2f) (Alarm: %t)", c.User, c.System, c.Alarm)
}

func (Utils) calCpuUsage() {
	percent, _ := cpu.Percent(time.Second, true)
	cpuInfo = CpuInfo{percent[cpu.CPUser], percent[cpu.CPSys], systemAlarm}

	if percent[cpu.CPSys] >= float64(systemAlarmLimit) && percent[cpu.CPUser] >= float64(systemAlarmLimit) {
		systemAlarm = true
	} else {
		systemAlarm = false
	}

	fmt.Printf("%v\n", cpuInfo)

	// fmt.Println("====== start cpu usage =====")
	// fmt.Printf("  User: %.2f\n", percent[cpu.CPUser])
	// fmt.Printf("  Nice: %.2f\n", percent[cpu.CPNice])
	// fmt.Printf("   Sys: %.2f\n", percent[cpu.CPSys])
	// fmt.Printf("  Intr: %.2f\n", percent[cpu.CPIntr])
	// fmt.Printf("  Idle: %.2f\n", percent[cpu.CPIdle])
	// fmt.Printf("States: %.2f\n", percent[cpu.CPUStates])
	// fmt.Printf("System Alarm: %t\n", systemAlarm)
	// fmt.Println("\n====== end cpu usage =====\n")
}

func (Utils) showStatus(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Total Connections:", totalConnections)
	// fmt.Fprintln(w, "Remain Connections:", currentConnections)
	// fmt.Fprintln(w, "Processed Rate:", processCounter.Rate(), "requests/seconds")
	// fmt.Fprintln(w, "Current Request Rate", requestCounter.Rate(), "requests/seconds")
	// fmt.Fprintf(w, "systemAlarm: %t", systemAlarm)

	info := Statistic{
		TotalConnections:   totalConnections,
		CurrentConnections: currentConnections,
		FailConnections:    failConnections,
		RequestRate:        requestCounter.Rate(),
		ProcessJob:         (totalConnections - currentConnections - failConnections),
		ProcessRate:        processCounter.Rate(),
		SystemAlarm:        systemAlarm,
	}

	data, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (Utils) lauchHttpServer() {
	http.HandleFunc("/status", utils.showStatus)
	go http.ListenAndServe(":8080", nil)
}

func (Utils) fechRemote() (bool, error) {

	resp, err := http.Get("http://" + endPointUrl + "/")

	if err != nil {
		fmt.Println("remote site can't connect!")
		fmt.Println("")
		return false, errors.New("remote site can't connect!")
	}

	log.Println("Status: ", resp.StatusCode)
	if resp.StatusCode == 200 {
		processRate = processCounter.Rate()
		if currentConnections > 0 {
			atomic.AddInt32(&currentConnections, -1)
		}
		log.Println("\n\n====>\nTotal Connections: ", totalConnections, "\nRemain Jobs: ", currentConnections, "\nProcessed Jobs", (totalConnections - currentConnections - failConnections), "\n<====\n\n")
		log.Println(" PROCESS RATE:", processRate, "jobs/second\n\n")
		fmt.Println("\nNomal Access")
	} else {
		fmt.Println("\nForbidden Access")
	}
	return true, nil
}
