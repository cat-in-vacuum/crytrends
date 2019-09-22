package main

import (
	"context"
	"fmt"
	"github.com/cat-in-vacuum/crytrade/providers/cryptowat"
	"net/http"
	"time"
)

func main() {
	data := cryptowat.AssetsContainer{}
	data.SetAsset("kraken", "btcusd")
	data.SetAsset("bitfinex", "ltcbtc")

	cw := cryptowat.New(&http.Client{}, &data, cryptowat.RetryPolicy{})
	ctx := cryptowat.Context{
		Context: context.Background(),
	}
	params := cryptowat.OHLCParams{
		Periods: []string{"60", "180"},
	}
	ts := time.Now()
	resp, err := cw.GetAllOHLC(ctx, params)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(time.Since(ts))
	fmt.Println(resp)
}
