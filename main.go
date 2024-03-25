package main

import "flag"

var (
	Scan    string = "stdin"
	Limit   int    = 10000
	Ports   string = "80,8000"
	Timeout int    = 5
	Verbose bool   = true
)

func init() {
	flag.StringVar(&Scan, "scan", "stdin", "scan type (stdin, cidr, file)")
	flag.IntVar(&Limit, "limit", 10000, "Thread limit")
	flag.StringVar(&Ports, "ports", "80,8000", "Ports to scan (Splited by comma)")
	flag.IntVar(&Timeout, "timeout", 5, "Timeout (seconds)")
	flag.BoolVar(&Verbose, "verbose", true, "verbose")

	flag.Parse()

	// TODO: Validate
}

func main() {

}
