package main

import (
	"context"
	"fmt"
	"github.com/cat-in-vacuum/crytrade/db"
	"github.com/cat-in-vacuum/crytrade/providers"
	"github.com/cat-in-vacuum/crytrade/providers/cryptowat"
	"github.com/cat-in-vacuum/crytrade/service"
	"net/http"
)

func main() {
	data := cryptowat.AssetsContainer{}
	data.SetAsset("kraken", "btcusd")
	data.SetAsset("bitfinex", "ltcbtc")

	cw := cryptowat.New(&http.Client{}, &data, cryptowat.RetryPolicy{})
	ctx := providers.Context{
		Context: context.Background(),
	}
	params := cryptowat.OHLCParams{
		Periods: []string{"60", "180"},
	}

	storage, _ := db.NewPostgres(fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", "admin" , "admin" , "postgres"))

    db.SetRepository(storage)

	srvs := service.NewService(storage, cw)
	srvs.StartStoringWithInterval(ctx, data, params)
}
