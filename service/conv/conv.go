package conv

import (
	"github.com/cat-in-vacuum/crytrade/db/cryptowat/schema"
	"github.com/cat-in-vacuum/crytrade/providers/cryptowat/types"
)

// пакет, ответсвенный за препроцессинг типов ответов поставщика в типы бд
// по идее, для каждого поставщика данных надо реализовать свой под-пакет
// с конвертациями

func OHLCRespToModels(period string, market, pair string, OHLC *types.OHLCResp) (out *schema.OHLCSchema) {
	 return &schema.OHLCSchema{
		QueryTime: OHLC.QueryTime,
		Market: market,
		Pair: pair,
		Period: period,
		OHLC: OHLC.Result[period],
	}
}
