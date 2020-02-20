/*
A very simple TCP server written in Go.
*/
package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	addr = ""
	port = 8000
)

var (
	listenAddr         string
	currentConnections int32
	totalConnections   int32
	failConnections    int32
	fillInterval       = time.Millisecond * 33 // let leaky bucket overflow rate to 33ms
	ticker             = time.NewTicker(fillInterval)
	requestCounter     = NewRateCounter(1 * time.Second)
	processCounter     = NewRateCounter(1 * time.Second)
	endPointUrl        = "localhost:3000"
	requestRate        int64
	processRate        int64
	cpuPercent         float64
	systemAlarm        = false
	systemAlarmLimit   = 90
	cpuInfo            CpuInfo
	rmutex             = &sync.Mutex{}
	utils              = Utils{}
)

func main() {

	src := addr + ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", src)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Println("TCP server start and listening on", src)
	totalConnections = 0
	failConnections = 0
	currentConnections = 0

	utils.lauchHttpServer()

	go func() {
		for {
			utils.calCpuUsage()
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Some connection error: %s", err)
			log.Fatal(err)
		}

		log.Println("\n\n==== Got a Request ==== ")
		atomic.AddInt32(&totalConnections, 1)
		requestCounter.Incr(1)
		requestRate = requestCounter.Rate()
		log.Println(" Current Request Rate:", requestRate, "requests/second\n\n")

		if !systemAlarm {
			atomic.AddInt32(&currentConnections, 1)
			go handleConnection(conn)
		} else {
			atomic.AddInt32(&failConnections, 1)
			fmt.Println("=== 503 Service Unavaliable ===")
			conn.Close()
		}
	}
}

func handleConnection(conn net.Conn) (bool, error) {
	// Close the connection when you're done with it.
	defer conn.Close()
	// Give the connection 1 minutes live-time
	err := conn.SetDeadline(time.Now().Add(1 * time.Minute))
	if err != nil {
		log.Fatalln("CONN TIMEOUT")
	}
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from: " + remoteAddr)
	fmt.Println("Current Connections: ", currentConnections)

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	for {
		// Read the incoming connection into the buffer.
		reqLen, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Disconned from ", remoteAddr)
				break
			} else {
				fmt.Println("Error reading:", err.Error())
				break
			}
		}
		// Send a response back to person contacting us.
		conn.Write([]byte("Message received.\n"))
		msg := string(buf[:reqLen])
		select {
		case <-ticker.C:
			if msg == "quit" {
				fmt.Println("=== Do nothing ===")
				// log.Println("len: %d, recv: %s\n", reqLen, msg)
				break
			} else {
				// log.Println("len: %d, recv: %s\n", reqLen, msg)
				processCounter.Incr(1)
				if !systemAlarm {
					go utils.fechRemote()
				}
			}
		}
	}
	return true, nil
}
