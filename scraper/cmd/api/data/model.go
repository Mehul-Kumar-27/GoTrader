package data

type Stock struct {
	Name     string
	Ticker   string
	Exchange string
	Price    string
}

func NewStock(name string, ticker string, exchange string) *Stock {
	return &Stock{
		Name:     name,
		Ticker:   ticker,
		Exchange: exchange,
	}
}
