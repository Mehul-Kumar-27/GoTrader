package main

import (
	"fmt"
	"gotrader/scraper/cmd/api/data"
	"sync"

	"github.com/gocolly/colly"
)

var url = "https://www.google.com/finance/quote/"

func main() {
	list := data.IndianStocks
	stocks := make([]data.Stock, len(list))
	priceCh := make(chan data.Stock, 4)

	i := 0
	for name, ticker := range list {
		st := data.NewStock(name, ticker, "NSE")
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

func aggregateDataFromGoroutine(priceCh chan data.Stock) {

	for stkResult := range priceCh {
		fmt.Printf("Received Stock: %s, Exchange: %s, Price: %s\n", stkResult.Ticker, stkResult.Exchange, stkResult.Price)
	}
}

func getPriceOfStock(stockList []data.Stock, priceCh chan data.Stock, wg *sync.WaitGroup, routine int) {
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
