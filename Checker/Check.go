package Checker

import (
	"ProxyJuice/Core"
	"ProxyJuice/Utility"
	"bytes"
	"net"
	"time"
)

type Proxy struct {
	Address string
	Port    int
	Joined  string
}

func Check(address string, port int) {
	proxy := &Proxy{Address: address, Port: port, Joined: address + ":" + Utility.FastIntToStr(port)}

	if checkForSocks5(proxy) {
		return
	} else if checkForSocks4(proxy) {
		return
	} else if checkForHttp(proxy) {
		return
	}

	//Core.LogChannel <- &Core.JuicyLog{
	//	Message: fmt.Sprintf("Couldn't Detect %s As Proxy", proxy.Joined),
	//	Level:   Core.OK,
	//	Sender:  "Checker.Check",
	//}

	proxy = nil
}

func checkForHttp(proxy *Proxy) bool {
	conn, err := net.DialTimeout("tcp", proxy.Joined, Core.Timeout)
	if err != nil {
		return false
	}

	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(Core.Timeout))
	if err != nil {
		return false
	}

	_, err = conn.Write(Core.HTTPGreeting)
	if err != nil {
		return false
	}

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return false
	}

	if bytes.Contains(buffer, []byte("200 Connection established")) {
		Core.SuccessChannel <- &Core.Response{Type: Core.HTTP, Address: proxy.Address, Port: proxy.Port}
		return true
	} else if bytes.Contains(buffer, []byte("407 Proxy Authentication Required")) {
		return checkPasswordsHTTP(proxy)
	}

	return false
}

func checkForSocks5(proxy *Proxy) bool { // success = true
	conn, err := net.DialTimeout("tcp", proxy.Joined, Core.Timeout)
	if err != nil {
		return false
	}

	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(Core.Timeout))
	if err != nil {
		return false
	}

	_, err = conn.Write(Core.SOCKS5Greeting)
	if err != nil {
		return false
	}

	buffer := make([]byte, 2)
	_, err = conn.Read(buffer)
	if err != nil {
		return false
	}

	if buffer[0] != 0x05 {
		return false
	}

	if buffer[1] == 0x00 {
		Core.SuccessChannel <- &Core.Response{Type: Core.SOCKS5, Address: proxy.Address, Port: proxy.Port}
		return true
	} else if buffer[1] == 0x02 {
		return checkPasswordsSOCKS5(proxy)
	}

	return false
}

func checkForSocks4(proxy *Proxy) bool {
	conn, err := net.DialTimeout("tcp", proxy.Joined, Core.Timeout)
	if err != nil {
		return false
	}

	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(Core.Timeout))
	if err != nil {
		return false
	}

	_, err = conn.Write(Core.SOCKS4Greeting)
	if err != nil {
		return false
	}

	buffer := make([]byte, 8)
	_, err = conn.Read(buffer)
	if err != nil {
		return false
	}

	if buffer[0] == 0x00 && buffer[1] == 0x5A {
		Core.SuccessChannel <- &Core.Response{Type: Core.SOCKS4, Address: proxy.Address, Port: proxy.Port}
		return true
	}

	return false
}

func checkPasswordsSOCKS5(proxy *Proxy) bool {
	for _, username := range Core.Usernames {
		for _, password := range Core.Passwords {
			payload := Core.BuildLoginPayloadSOCKS5(username, password)

			conn, err := net.DialTimeout("tcp", proxy.Joined, Core.Timeout)
			if err != nil {
				return false
			}

			err = conn.SetDeadline(time.Now().Add(Core.Timeout))
			if err != nil {
				_ = conn.Close()
				return false
			}

			_, err = conn.Write(payload)
			if err != nil {
				_ = conn.Close()
				return false
			}

			buffer := make([]byte, 2)
			_, err = conn.Read(buffer)
			if err != nil {
				_ = conn.Close()
				return false
			}

			if buffer[0] == 0x01 && buffer[1] == 0x00 {
				Core.SuccessChannel <- &Core.Response{Type: Core.SOCKS5, Address: proxy.Address, Port: proxy.Port, Username: username, Password: password}
				_ = conn.Close()
				return true
			}
		}
	}

	return false
}

func checkPasswordsHTTP(proxy *Proxy) bool {
	for _, username := range Core.Usernames {
		for _, password := range Core.Passwords {
			payload := Core.BuildLoginPayloadHTTP(username, password)

			conn, err := net.DialTimeout("tcp", proxy.Joined, Core.Timeout)
			if err != nil {
				return false
			}

			err = conn.SetDeadline(time.Now().Add(Core.Timeout))
			if err != nil {
				_ = conn.Close()
				return false
			}

			_, err = conn.Write(payload)
			if err != nil {
				_ = conn.Close()
				return false
			}

			buffer := make([]byte, 1024)
			_, err = conn.Read(buffer)
			if err != nil {
				_ = conn.Close()
				return false
			}

			if bytes.Contains(buffer, []byte("200 Connection established")) {
				_ = conn.Close()
				Core.SuccessChannel <- &Core.Response{Type: Core.HTTP, Address: proxy.Address, Port: proxy.Port}
				return true
			}
		}
	}

	return false
}
