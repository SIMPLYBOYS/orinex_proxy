package main

import (
	"fmt"
	"net"
	"time"
)

// FetchRemote ...
func FetchRemote(msg string) {
	res, err := sendTCP("127.0.0.1:8000", msg)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(res)
	}
}

// TestRemoteBursty ...
func TestRemoteBursty(count int, msg string) {
	//===========

	burstyLimiter := make(chan time.Time, count)

	for i := 0; i < count; i++ {
		burstyLimiter <- time.Now()
	}

	burstyRequests := make(chan int, count)
	for i := 1; i <= count; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)
	for req := range burstyRequests {
		<-burstyLimiter
		FetchRemote(msg)
		fmt.Println("request", req, time.Now())
	}
}

func sendTCP(addr, msg string) (string, error) {
	// connect to this socket
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// send to socket
	conn.Write([]byte(msg))

	// listen for reply
	bs := make([]byte, 1024)
	len, err := conn.Read(bs)
	if err != nil {
		return "", err
	} else {
		return string(bs[:len]), err
	}
}
