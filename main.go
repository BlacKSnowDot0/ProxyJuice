package main

import (
	"ProxyJuice/Checker"
	"ProxyJuice/Core"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	Core.LogChannel <- &Core.JuicyLog{
		Message: "Starting ProxyJuice",
		Level:   Core.OK,
		Sender:  "ProxyJuice",
	}

	Checker.Start()
}
