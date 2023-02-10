package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"strings"
	"syscall/js"
)

var document = js.Global().Get("document")

func getElementByID(id string) js.Value {
	return document.Call("getElementById", id)
}

func main() {

	quit := make(chan struct{}, 0)

	// See example 2: Enable the stop button
	stopButton := getElementByID("stop")
	stopButton.Set("disabled", false)
	stopButton.Set("onclick", js.FuncOf(func(js.Value, []js.Value) interface{} {
		println("stopping")
		stopButton.Set("disabled", true)
		quit <- struct{}{}
		return nil
	}))

	// Instantiate default collector
	c := colly.NewCollector()

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	// On every a element which has href attribute call callback
	c.OnResponse(func(r *colly.Response) {

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(r.Body)))
		if err != nil {
			return
		}
		pathid := "#hprt-table > tbody > tr.js-rt-block-row.e2e-hprt-table-row.hprt-table-cheapest-block.hprt-table-cheapest-block-fix.js-hprt-table-cheapest-block > td.hp-price-left-align.hprt-table-cell.hprt-table-cell-price > div > div > div:nth-child(1) > div:nth-child(2) > div > span"
		price := doc.Find(pathid)
		var children = price.Text()
		fmt.Println(children)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("sec-fetch-mode", "no-cors")
		r.Headers.Set("sec-fetch-dest", "empty")
		r.Headers.Set("sec-fetch-site", "cross-site")
	})

	url := "https://www.booking.com/hotel/no/spitsbergen.en-gb.html?no_rooms=1&checkin=2023-04-03&checkout=2023-04-07&group_adults=2&group_children=0&req_adults=2&req_children=0&dist=0&type=total&selected_currency=NOK"
	c.Visit(url)
	url = "https://www.booking.com/hotel/es/sol-pelicanos-ocas.es.html?no_rooms=1&checkin=2023-04-03&checkout=2023-04-07&group_adults=2&group_children=0&req_adults=2&req_children=0&dist=0&type=total&selected_currency=EUR"
	c.Visit(url)

	<-quit
}
