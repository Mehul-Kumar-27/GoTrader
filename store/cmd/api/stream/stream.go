package stream

import (
	lg "gotrader/logger"
	pb "gotrader/proto"
	"gotrader/store/cmd/api/consumer"
	"log"
)

var logger *log.Logger

func init() {
	logger = lg.CreateCustomLogger("stream")
}

type Stream struct {
	pb.UnimplementedStockServiceServer
}

func (s *Stream) GetStocks(req *pb.ExchangeRequest, stream pb.StockService_GetStocksServer) error {
	logger.Printf("Received request for exchange: %s", req.Exchange)
	if !consumer.NseChannelInitialized {
		consumer.SubscribleToNseChannel()
		for stock := range consumer.NseChn {
			if err := stream.Send(stock); err != nil {
				logger.Fatalf("Error occurred while sending the stock to the stream %s", err)
			}
		}
	}
	return nil
}
