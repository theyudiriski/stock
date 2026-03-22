package config

type Redis struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func LoadRedis() Redis {
	return Redis{
		Host:     GetString("redis.host"),
		Port:     GetString("redis.port"),
		Password: GetStringOptional("redis.password"),
		DB:       GetIntOptional("redis.db"),
	}
}
