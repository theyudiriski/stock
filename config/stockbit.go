package config

type Stockbit struct {
	Tokens  []string
	BaseURL string
}

func LoadStockbit() Stockbit {
	return Stockbit{
		Tokens:  GetStringSlice("stockbit.tokens"),
		BaseURL: GetString("stockbit.base_url"),
	}
}
