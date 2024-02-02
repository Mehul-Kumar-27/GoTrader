package consumer


type KafkaMessage struct {
	StockName     string
	StockExchange string
	StockTicker   string
	StockPrice    string
}