package consumer

import (
	"context"
	"encoding/json"
	pb "gotrader/proto"
	"gotrader/scraper/cmd/api/logger"
	"log"

	"github.com/segmentio/kafka-go"
)

var reader *kafka.Reader
var ConsumerLogger *log.Logger
var NseChn chan *pb.Stock
var NseChannelInitialized bool = false

func init() {
	reader, _ = NewKafkaReader([]string{"broker:9092"}, "stock", 0)
	ConsumerLogger = logger.CreateCustomLogger("consumer")
	ConsumerLogger.Println("Consumer is initialized")

}

func NewKafkaReader(brokers []string, topic string, partition int) (*kafka.Reader, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		MaxBytes:  10e6,
		Partition: partition,
	})
	reader.SetOffset(2)
	return reader, nil
}

func GetTheReader() *kafka.Reader {
	return reader
}

func ConsumeMessages() {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			reader.Close()
			ConsumerLogger.Fatalf("Error occurred while reading the message from the kafka %s", err)
		} else {
			var kafkaMessage KafkaMessage
			err := json.Unmarshal(m.Value, &kafkaMessage)
			if err != nil {
				ConsumerLogger.Fatalf("Error occurred while unmarshalling the stock %s", err)
			}
			if NseChannelInitialized {
				var stk = &pb.Stock{
					Name:     kafkaMessage.StockName,
					Ticker:   kafkaMessage.StockTicker,
					Exchange: kafkaMessage.StockExchange,
					Price:    kafkaMessage.StockPrice,
				}
				NseChn <- stk
			}

			ConsumerLogger.Printf("Recieved Stock from at the offset %d : %s", m.Offset, string(m.Value))
		}
	}
}

func SubscribleToNseChannel() {
	if NseChn == nil {
		NseChn = make(chan *pb.Stock)
		NseChannelInitialized = true
	} else {
		ConsumerLogger.Println("NSE Channel is already initialized")
	}
}
