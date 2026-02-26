package httpclient

import (
	"fmt"
	"net/http"
	"net/url"
	"stock/config"
)

func New() *http.Client {
	cfg := config.LoadProxy()
	if !cfg.Enabled {
		return &http.Client{}
	}

	proxyURL, err := build(cfg.Scheme, cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	if err != nil {
		return &http.Client{}
	}

	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
}

func build(scheme, host string, port int, username, password string) (*url.URL, error) {
	raw := fmt.Sprintf("%s://%s:%d", scheme, host, port)
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	u.User = url.UserPassword(username, password)
	return u, nil
}
