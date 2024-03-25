package Parser

type ProxyType string

const (
	SOCKS4 ProxyType = "SOCKS4"
	SOCKS5 ProxyType = "SOCKS5"
	HTTP   ProxyType = "HTTP"
)

type Response struct {
	Type     ProxyType
	Address  string
	Port     int
	Username string
	Password string
}
