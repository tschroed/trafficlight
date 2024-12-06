package main

import (
	"strconv"
	"time"
	"flag"

	"github.com/tschroed/trafficlight"
	"github.com/tschroed/trafficlight/k8090"
	"github.com/tschroed/trafficlight/lcus"
)

var dFlag = flag.String("d", "k8090", "Driver to use (k8090 or lcus)")

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var tl trafficlight.TrafficLight
	var err error
	flag.Parse()
	switch *dFlag {
	case "k8090":
		tl, err = k8090.New("/dev/ttyACM0")
	case "lcus":
		tl, err = lcus.New("/dev/ttyUSB0")
	}
	check(err)
	defer tl.Set(0)
	s, err := strconv.Atoi(flag.Args()[0])
	check(err)
	for _, arg := range flag.Args()[1:] {
		n, err := strconv.Atoi(arg)
		check(err)
		tl.Set(uint8(n))
		time.Sleep(time.Duration(s) * time.Millisecond)
	}
}
