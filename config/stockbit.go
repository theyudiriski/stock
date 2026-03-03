package config

type Stockbit struct {
	BaseURL       string
	Tokens        []string
	WebviewTokens []string
}

func LoadStockbit() Stockbit {
	return Stockbit{
		BaseURL:       GetString("stockbit.base_url"),
		Tokens:        GetStringSlice("stockbit.tokens"),
		WebviewTokens: GetStringSlice("stockbit.webview_tokens"),
	}
}
