package main

import (
	"fmt"
	"github.com/cat-in-vacuum/crytrade/providers/cryptowat"
	"net/http"
)

func main() {
	data := cryptowat.AssetsContainer{}
	data.SetAsset("kraken", "btcusd")
	cw := cryptowat.New(&http.Client{}, &data)
	params := cryptowat.OHLCParams{
		Periods: []string{"60", "180"},
	}
	resp, err := cw.GetAllOHLC(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}
