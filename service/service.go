package service

import (
	"github.com/cat-in-vacuum/crytrade/db"
	"github.com/cat-in-vacuum/crytrade/providers"
	"github.com/cat-in-vacuum/crytrade/providers/cryptowat"
	"github.com/cat-in-vacuum/crytrade/scheduler"
	"github.com/cat-in-vacuum/crytrade/service/conv"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// по идее, вместо cryptowat.Client тут должен быть слайс интерфейсов типа Provider
// которые реализуют метод поставки типов providers.OHLCAssets
type Service struct {
	db        db.Repository
	scheduler *scheduler.Scheduler
	client    *cryptowat.Client
}

func NewService(
	db db.Repository,
	client *cryptowat.Client,
	scheduler *scheduler.Scheduler,
) *Service {
	return &Service{db: db, scheduler: scheduler, client: client}
}

func (s Service) StartStoringOHLC(ctx providers.Context, assets cryptowat.AssetsContainer, params cryptowat.OHLCParams, ) error {
	ctx.SetID()
	resp, err := s.client.GetAllOHLC(ctx, params)
	if err != nil {
		return errors.Wrap(err, "error in Service.StartStoringOHLC")
	}

	for market, asset := range resp {
		for pair, ohlc := range asset {
			log.Debug().Str("pair", pair).Str("market", market).Msgf("start storing asset in service.StartStoringOHLC(); ReqID:%s", ctx.GetID())
			for period, _ := range ohlc.Result {
				err := s.db.InsertOHLC(conv.OHLCRespToModels(period, market, pair, &ohlc))

				if err != nil {
					log.Error().Err(err).Str("pair", pair).Str("market", market).Msgf("Asset not saved, ReqID:%s", ctx.GetID())
					continue
				}
			}
		}
	}

	return nil
}
