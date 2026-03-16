package config

type Telegram struct {
	Bot TelegramBot
}

type TelegramBot struct {
	Token string
}

func LoadTelegram() Telegram {
	return Telegram{
		Bot: TelegramBot{
			Token: GetString("telegram.bot.token"),
		},
	}
}
