package config

type Stockbit struct {
	BaseURL       string
	Tokens        map[string]string
	WebviewTokens []string
}

func LoadStockbit() Stockbit {
	return Stockbit{
		BaseURL:       GetString("stockbit.base_url"),
		Tokens:        GetStringMapString("stockbit.tokens"),
		WebviewTokens: GetStringSlice("stockbit.webview_tokens"),
	}
}
