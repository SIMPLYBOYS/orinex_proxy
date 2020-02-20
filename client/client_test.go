package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestRemoteQuery(t *testing.T) {

	var (
		burstyNum int
		msg       string
	)

	for _, arg := range flag.Args() {
		arglist := strings.Split(arg, "=")

		if len(arglist) != 4 {
			arglist = strings.Split(arg, " ")
		}

		if len(arglist) != 4 {
			fmt.Printf(" unknown arg:%v\n", arg)
			continue
		}
		arg0 := strings.TrimSpace(arglist[0])
		if arg0 == "c" || arg0 == "C" {
			burstyNum, _ = strconv.Atoi(strings.TrimSpace(arglist[1]))
		} else {
			fmt.Printf(" unknown format:%v\n", arg)
			continue
		}

		arg1 := strings.TrimSpace(arglist[2])
		if arg1 == "m" || arg1 == "M" {
			msg = strings.TrimSpace(arglist[3])
		} else {
			fmt.Printf(" unknown format:%v\n", arg)
			continue
		}
	}
	TestRemoteBursty(burstyNum, msg)
}

func BenchmarkRemoteQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestRemoteBursty(1, "hi")
	}
}
