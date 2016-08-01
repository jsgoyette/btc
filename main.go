// why only this comment?
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// CoinBaseResponse - The response from Coinbase
type CoinBaseResponse struct {
	Amount string `json:"amount"`
}

// BittrexResponse - The response from Bittrex
type BittrexResponse struct {
	Result []BittrexMarket `json:"result"`
}

// BittrexMarket - The returned markets from Bittrex
type BittrexMarket struct {
	MarketName string
	Bid        float32
	Ask        float32
}

func handleRequest(url string) ([]byte, error) {

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func makeBittrexRequest(url string, res chan BittrexResponse) {

	contents, _ := handleRequest(url)
	btxres := new(BittrexResponse)

	if err := json.Unmarshal([]byte(contents), btxres); err != nil {
		fmt.Printf("%v", err)
	}

	res <- *btxres
}

func makeCoinBaseRequest(url string, res chan CoinBaseResponse) {

	contents, _ := handleRequest(url)
	cbres := new(CoinBaseResponse)

	if err := json.Unmarshal([]byte(contents), cbres); err != nil {
		fmt.Printf("%v", err)
	}

	res <- *cbres
}

var acceptedMarkets = map[string]bool{
	"BTC-BLK":  true,
	"BTC-DASH": true,
	"BTC-DGD":  true,
	"BTC-DOGE": true,
	"BTC-ETC":  true,
	"BTC-ETH":  true,
	"BTC-LTC":  true,
}

var yellow = "\033[1;33m"
var nocolor = "\033[0m"

func main() {

	buyRes := make(chan CoinBaseResponse)
	sellRes := make(chan CoinBaseResponse)
	bittrexRes := make(chan BittrexResponse)

	go makeCoinBaseRequest("https://api.coinbase.com/v1/prices/buy", buyRes)
	go makeCoinBaseRequest("https://api.coinbase.com/v1/prices/sell", sellRes)
	go makeBittrexRequest("https://bittrex.com/api/v1.1/public/getmarketsummaries", bittrexRes)

	buy := <-buyRes
	sell := <-sellRes
	bittrex := <-bittrexRes

	fmt.Println(yellow + "BTC" + nocolor)
	fmt.Printf("Bid:\t$%s\nAsk:\t$%s\n", sell.Amount, buy.Amount)

	for _, market := range bittrex.Result {
		if acceptedMarkets[market.MarketName] {
			fmt.Println(yellow + market.MarketName + nocolor)
			fmt.Printf("Bid:\t%.8f\nAsk:\t%.8f\n", market.Bid, market.Ask)
		}
	}
}
