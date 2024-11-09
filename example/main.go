package main

import (
	"os"
	"strconv"
	"time"

	"ryanairship.com/trafficlight/k8090"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	k, err := k8090.New("/dev/ttyACM0")
	check(err)
	defer k.Set(0)
	s, err := strconv.Atoi(os.Args[1])
	check(err)
	for _, arg := range os.Args[2:] {
		n, err := strconv.Atoi(arg)
		check(err)
		k.Set(uint8(n))
		time.Sleep(time.Duration(s) * time.Second)
	}
}
