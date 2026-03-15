package service

type TradeBookResponse struct {
	Data    TradeBook `json:"data"`
	Message string    `json:"message"`
}

type TradeBook struct {
	Book []TradeBookBook `json:"book"`
	Date string          `json:"date"`
}

type TradeBookBook struct {
	Time string        `json:"time"`
	Buy  TradeBookItem `json:"buy"`
	Sell TradeBookItem `json:"sell"`
}

type TradeBookItem struct {
	Lot       string `json:"lot"`
	Frequency string `json:"frequency"`
}
