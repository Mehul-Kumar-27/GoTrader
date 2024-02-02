package consumer

import (
	"context"
	pb "gotrader/proto"
	"gotrader/scraper/cmd/api/logger"
	"log"

	"github.com/segmentio/kafka-go"
)

var reader *kafka.Reader
var ConsumerLogger *log.Logger
var nseChn chan *pb.Stock

func init() {
	reader, _ = NewKafkaReader([]string{"localhost:19092"}, "stock", 0)
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
			ConsumerLogger.Printf("Recieved Stock from at the offset %d : %s", m.Offset, string(m.Value))
		}
	}
}
