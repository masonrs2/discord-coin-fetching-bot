package structs

type Token struct {
	Id     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type Response struct {
	MarketData struct {
		CurrentPrice struct {
			Usd float64 `json:"usd"`
		} `json:"current_price"`
	} `json:"market_data"`
}