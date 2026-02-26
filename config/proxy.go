package config

type Proxy struct {
	Enabled  bool
	Scheme   string
	Host     string
	Port     int
	Username string
	Password string
}

func LoadProxy() Proxy {
	return Proxy{
		Enabled:  GetBool("proxy.enabled"),
		Scheme:   GetString("proxy.scheme"),
		Host:     GetString("proxy.host"),
		Port:     GetInt("proxy.port"),
		Username: GetString("proxy.username"),
		Password: GetString("proxy.password"),
	}
}
