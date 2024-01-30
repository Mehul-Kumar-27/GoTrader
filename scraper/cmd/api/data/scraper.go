package data

import (
	"fmt"
	"gotrader/scraper/cmd/api/producer"
	"sync"

	"github.com/gocolly/colly"
)

var url = "https://www.google.com/finance/quote/"

func StartScraping() {
	list := IndianStocks
	stocks := make([]Stock, len(list))
	priceCh := make(chan Stock, 4)

	i := 0
	for name, ticker := range list {
		st := NewStock(name, ticker, "NSE")
		stocks[i] = *st
		i++
	}

	var wg sync.WaitGroup
	var sectionSize = len(stocks) / 4
	for i := 0; i < 4; i++ {
		startIndex := i * sectionSize
		endIndex := (i + 1) * sectionSize

		if endIndex > len(stocks) {
			endIndex = len(stocks)
		}

		wg.Add(1)
		go getPriceOfStock(stocks[startIndex:endIndex], priceCh, &wg, i)
	}

	go aggregateDataFromGoroutine(priceCh)

	wg.Wait()
	close(priceCh)
}

func aggregateDataFromGoroutine(priceCh chan Stock) {

	for stkResult := range priceCh {
		producer.PublishStockToKafka(stkResult.Name, stkResult.Ticker, stkResult.Exchange, stkResult.Price, []string{"localhost:19092"})
	}
}

func getPriceOfStock(stockList []Stock, priceCh chan Stock, wg *sync.WaitGroup, routine int) {
	defer wg.Done()

	for _, stock := range stockList {
		c := colly.NewCollector()
		c.OnHTML("div.YMlKec.fxKbKc", func(e *colly.HTMLElement) {
			stock.Price = e.Text
			priceCh <- stock
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Printf("Request URL: %s failed with response: %v\nError: %v\n", r.Request.URL, r, err)
		})

		c.Visit(url + stock.Ticker + ":" + stock.Exchange)
		c.Wait()
	}

}
