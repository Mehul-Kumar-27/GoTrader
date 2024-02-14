package listner

import (
	"context"
	lg "gotrader/logger"
	pb "gotrader/proto"
	"gotrader/server/cmd/api/routes"
	"log"
)

var logger *log.Logger

func init() {
	logger = lg.CreateCustomLogger("listner")
	NseChannelManager = make([]*routes.ClientHandeller, 0)

}

func AddClient(client *routes.ClientHandeller) {
	NseChannelManager = append(NseChannelManager, client)
}

func RemoveClient(client *routes.ClientHandeller) {
	for i, c := range NseChannelManager {
		if c == client {
			NseChannelManager = append(NseChannelManager[:i], NseChannelManager[i+1:]...)
			break
		}
	}
}

var NseChannelManager []*routes.ClientHandeller

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
		logger.Printf("Stock: %v price: %v", message.Name, message.Price)

		for _, client := range NseChannelManager {
			client.SendToNseChan(message)
		}

	}

}
