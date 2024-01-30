package producer

import (
	"context"
	"encoding/json"

	"gotrader/scraper/cmd/api/logger"

	"github.com/segmentio/kafka-go"
)

type KafkaMessage struct {
	StockName     string
	StockExchange string
	StockTicker   string
	StockPrice    string
}

var writer *kafka.Writer
var i int = 0

func init() {
	writer = CreateCustomWriter([]string{"localhost:19092"}, "stock")
}

func GetTheWriter() *kafka.Writer {
	return writer
}

func CreateCustomWriter(kafkaBrokers []string, topic string) *kafka.Writer {
	logger := logger.CreateCustomLogger("CreateCustomWriter")
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				logger.Fatalf("Error occured while writing to the kafka message from the completion section %s", err)
			} else {
				logger.Println(messages)
			}
		},
		AllowAutoTopicCreation: true,
	}
}

func PublishStockToKafka(name, ticker, exchange, price string, kafkaBrokers []string) {
	logger := logger.CreateCustomLogger("PublishStockToKafka")
	data, err := json.Marshal(KafkaMessage{
		StockName:     name,
		StockExchange: exchange,
		StockTicker:   ticker,
		StockPrice:    price,
	})

	if err != nil {
		logger.Fatalf("Error occurred while marshalling the stock %s", err)
	}
	logger.Println("Attempting to write to the kafka")
	er := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(name),
		Value: data,
	})

	if er != nil {
		logger.Fatalf("Error occurred while writing to the kafka %s", er)
	} else {
		i++
		logger.Printf("Total %v stocks published", i)
	}
}
