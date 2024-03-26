package Core

import (
	"ProxyJuice/Utility"
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"time"
)

var (
	Scan       = "stdin"
	Limit      = 10000
	Ports      = "80,8000"
	Timeout    time.Duration
	Usernames  = []string{"Misaka"}
	Passwords  = []string{"Misaka"}
	OutputFile = "output.txt"
	Verbose    = true
)

var (
	SuccessChannel = make(chan *Response)
	LogChannel     = make(chan *JuicyLog)
)

type ProxyType string
type JuiceLevel string

const (
	SOCKS4 ProxyType = "SOCKS4"
	SOCKS5 ProxyType = "SOCKS5"
	HTTP   ProxyType = "HTTP"
)

const (
	SUCCESS JuiceLevel = "SUCCESS"
	OK      JuiceLevel = "OK"
	WARN    JuiceLevel = "WARN"
	BAD     JuiceLevel = "ERRBAD"
	DIE     JuiceLevel = "ERRDIE"
)

type Response struct {
	Type     ProxyType
	Address  string
	Port     int
	Username string
	Password string
}

type JuicyLog struct {
	Message string
	Level   JuiceLevel
	Sender  string
}

func init() {
	flag.StringVar(&Scan, "scan", "stdin", "scan type (stdin, cidr, file)")
	flag.IntVar(&Limit, "limit", 10000, "Thread limit")
	flag.StringVar(&Ports, "ports", "80,8000", "Ports to scan (Splited by comma)")
	timeout := flag.Int("timeout", 5, "Timeout (seconds)")
	usernames := flag.String("usernames", "usernames.txt", "Usernames file/Url")
	passwords := flag.String("passwords", "passwords.txt", "Passwords file/Url")
	flag.StringVar(&OutputFile, "output", "output.txt", "Output file")
	flag.BoolVar(&Verbose, "verbose", true, "verbose")

	flag.Parse()

	if Scan != "stdin" && Scan != "cidr" && Scan != "file" {
		log.Fatal("Invalid scan type. Please use stdin, cidr, or file.")
	}

	Timeout = time.Duration(*timeout) * time.Second
	Usernames = handleUP(*usernames)
	Passwords = handleUP(*passwords)

	if len(Usernames) == 0 {
		Usernames = []string{"Misaka"}
	}

	if len(Passwords) == 0 {
		Passwords = []string{"Misaka"}
	}

	go SuccessChannelHandler()
	go LogChannelHandler()
}

func handleUP(input string) []string {
	if len(input) > 3 && input[:4] == "http" {
		return Utility.GatherLines(input)
	}

	if _, err := os.Stat(input); err == nil {
		file, err := os.Open(input)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		var lines []string

		for scanner.Scan() {
			lines = append(lines, strings.TrimSpace(scanner.Text()))
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		return lines
	}

	return []string{}
}
