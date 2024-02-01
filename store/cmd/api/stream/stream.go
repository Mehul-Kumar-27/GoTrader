package stream

import (
	lg "gotrader/logger"
	pb "gotrader/proto"
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
	for i := 0; i < 10; i++ {
		stock := &pb.Stock{
			Name:  string(rune(i)),
			Price: "100",
		}
		if err := stream.Send(stock); err != nil {
			logger.Printf("Error sending stock: %v", err)
		}
	}
	return nil
}
