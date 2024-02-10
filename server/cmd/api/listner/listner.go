package listner

import (
	"context"
	lg "gotrader/logger"
	pb "gotrader/proto"
	"log"
)

var logger *log.Logger

func init() {
	logger = lg.CreateCustomLogger("listner")
}

func GetStocksForNse(client pb.StockServiceClient) {
	stream, err := client.GetStocks(context.Background(), &pb.ExchangeRequest{Exchange: "NSE"})
	if err != nil {
		logger.Fatalf("Error getting stocks: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err != nil {
			logger.Fatalf("Error receiving message: %v", err)
		}
		logger.Printf("Received stock: %v price: %v", message.Name, message.Price)
	}

}
