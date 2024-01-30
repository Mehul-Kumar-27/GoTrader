package main

import (
	"gotrader/scraper/cmd/api/data"
	"gotrader/scraper/cmd/api/logger"
	"time"
)

func main() {
	logger := logger.CreateCustomLogger("main")
	for {
		data.StartScraping()
		logger.Println("Sleeping for 30 seconds")
		time.Sleep(30 * time.Second)
	}
}
