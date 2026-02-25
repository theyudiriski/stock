package config

import "time"

type Database struct {
	Leader DatabaseConfig
}

type DatabaseConfig struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DB              string
	Scheme          string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	SSLMode         string
}

func LoadDatabase() Database {
	return Database{
		Leader: DatabaseConfig{
			Host:            GetString("database.host"),
			Port:            GetString("database.port"),
			Username:        GetString("database.username"),
			Password:        GetStringOptional("database.password"),
			DB:              GetString("database.db"),
			Scheme:          GetString("database.scheme"),
			MaxIdleConns:    GetIntOptional("database.max_idle_conns"),
			MaxOpenConns:    GetIntOptional("database.max_open_conns"),
			ConnMaxLifetime: GetDurationOptional("database.conn_max_lifetime"),
			SSLMode:         GetStringOptional("database.ssl_mode"),
		},
	}
}
