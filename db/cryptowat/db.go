package cryptowat

import (
	"context"
	"github.com/cat-in-vacuum/crytrade/db/cryptowat/schema"
)

type Repository interface {
	Close()
	InsertOHLC(ctx context.Context, ohlc schema.OHLCSchema) error
	ListOHLC(ctx context.Context, skip uint64, take uint64) ([]schema.OHLCSchema, error)
}

// оставляем возможность переключать дб на лету
var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertOHLC(ctx context.Context, ohlc OHLCSchema) error {
	return impl.InsertOHLC(ctx, ohlc)
}

func ListOHLC(ctx context.Context, skip uint64, take uint64) ([]OHLCSchema, error) {
	return impl.ListOHLC(ctx, skip, take)
}

