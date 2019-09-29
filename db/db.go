package db

import (
	"github.com/cat-in-vacuum/crytrade/db/cryptowat/schema"
)

type Repository interface {
	Close()
	InsertOHLC(ohlc *schema.OHLCSchema) error
	ListOHLC(skip uint64, take uint64) ([]schema.OHLCSchema, error)
}

// оставляем возможность переключать дб на лету
var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertOHLC(ohlc *schema.OHLCSchema) error {
	return impl.InsertOHLC(ohlc)
}

func ListOHLC(skip uint64, take uint64) ([]schema.OHLCSchema, error) {
	return impl.ListOHLC(skip, take)
}

