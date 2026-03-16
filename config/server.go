package config

import "time"

type Server struct {
	Addr    string
	Timeout ServerTimeout
}

type ServerTimeout struct {
	Read  time.Duration
	Write time.Duration
	Idle  time.Duration
}

func LoadServer() Server {
	return Server{
		Addr: GetString("server.addr"),
		Timeout: ServerTimeout{
			Read:  GetDuration("server.timeout.read"),
			Write: GetDuration("server.timeout.write"),
			Idle:  GetDuration("server.timeout.idle"),
		},
	}
}
