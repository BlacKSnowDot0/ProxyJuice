package Core

import (
	"bytes"
	"encoding/base64"
)

var (
	SOCKS5Greeting = []byte{0x05, 0x02, 0x00, 0x02}
	SOCKS4Greeting = []byte{0x04, 0x01, 0x00, 0x35, 0x08, 0x08, 0x08, 0x08, 0x00, 0x00}
	HTTPGreeting   = []byte("CONNECT 1.1.1.1:80 HTTP/1.0\r\nHost: 1.1.1.1\r\n\r\n")
)

func BuildLoginPayloadSOCKS5(username, password string) []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x01)
	buf.WriteByte(byte(len(username)))
	buf.WriteString(username)
	buf.WriteByte(byte(len(password)))
	buf.WriteString(password)
	return buf.Bytes()
}

func BuildLoginPayloadHTTP(username string, password string) []byte {
	b64 := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	var buf bytes.Buffer
	buf.WriteString("CONNECT 1.1.1.1:80 HTTP/1.0\r\nHost: 1.1.1.1\r\n")
	buf.WriteString("Proxy-Authorization: Basic ")
	buf.WriteString(b64)
	buf.WriteString("\r\n\r\n")

	return buf.Bytes()
}
