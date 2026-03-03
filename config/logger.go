package config

type Logger struct {
	Level  string
	Format string
}

func LoadLogger() Logger {
	return Logger{
		Level:  GetString("logger.level"),
		Format: GetString("logger.format"),
	}
}
