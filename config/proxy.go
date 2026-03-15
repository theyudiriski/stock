package config

import "stock/internal/service"

type Proxy struct {
	Enabled  ProxyEnabled
	Scheme   string
	Host     string
	Port     int
	Username string
	Password string
}

type ProxyEnabled struct {
	Stockbit bool
}

func LoadProxy() Proxy {
	return Proxy{
		Enabled: ProxyEnabled{
			Stockbit: GetBool("proxy.enabled.stockbit"),
		},
		Scheme:   GetString("proxy.scheme"),
		Host:     GetString("proxy.host"),
		Port:     GetInt("proxy.port"),
		Username: GetString("proxy.username"),
		Password: GetString("proxy.password"),
	}
}

func (p *Proxy) IsEnabled(serviceName service.ServiceName) bool {
	switch serviceName {
	case service.ServiceNameStockbit:
		return p.Enabled.Stockbit
	default:
		return false
	}
}
