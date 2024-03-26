package Checker

import (
	"ProxyJuice/CIDRManager"
	"ProxyJuice/Core"
	"ProxyJuice/Utility"
	"ProxyJuice/pool"
	"bufio"
	"errors"
	"math/rand"
	"os"
	"strings"
)

type Checker struct {
	InputChannel chan string
	Pool         *pool.WorkerPool
}

func Start() {
	checker := Checker{
		InputChannel: make(chan string),
		Pool:         pool.New(Core.Limit),
	}

	go checker.Listen()

	switch Core.Scan {
	case "stdin":
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			checker.InputChannel <- strings.TrimSpace(scanner.Text())
		}

		if scanner.Err() != nil {
			Core.LogChannel <- &Core.JuicyLog{
				Message: scanner.Err().Error(),
				Level:   Core.DIE,
				Sender:  "Checker.Start",
			}
		}

	default:
		var lines []string
		if Utility.FileExists(Core.Scan) {
			lines = Utility.ReadFileLines(Core.Scan)
		}

		var cidrs []*CIDRManager.CIDRManager

		for _, line := range lines {
			cidrs = append(cidrs, CIDRManager.NewCIDRManager(line))
		}

		Amounts := len(cidrs)
		for Amounts > 0 {
			numb := rand.Intn(Amounts)
			cidr := cidrs[numb]

			ip, err := cidr.GetUnusedIP()
			if err != nil {
				if errors.Is(err, CIDRManager.EOCIDR) {
					cidrs = append(cidrs[:numb], cidrs[numb+1:]...)
					Amounts--
					continue
				}

				Core.LogChannel <- &Core.JuicyLog{
					Message: err.Error(),
					Level:   Core.BAD,
					Sender:  "Checker.Start",
				}
			}

			checker.InputChannel <- ip
		}
	}

	Core.LogChannel <- &Core.JuicyLog{
		Message: "Finished",
		Level:   Core.SUCCESS,
		Sender:  "Checker.Start",
	}
}

func (checker *Checker) Listen() {
	checker.Pool.Start()

	for data := range checker.InputChannel {
		for _, port := range Core.Ports {
			checker.Pool.Submit(func() {
				Check(data, port)
			})
		}
	}

	checker.Pool.Stop()
}
