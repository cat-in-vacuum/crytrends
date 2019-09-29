package cryptowat

import (
	"context"
	"database/sql"

	"github.com/cat-in-vacuum/crytrade/db/cryptowat/schema"
	_ "github.com/lib/pq"
)

const (
	queryInsert = "insert into CRYPTOWAT_OHLC_DATA(ID, QUERY_TIME, MARKET, PAIR, PERIOD, OHLC) values($1, $2, $3, $4, $5, $6)"
	queryList = "select * from CRYPTOWAT_OHLC_DATA order by ID desc offset $1 limit $2"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{
		db,
	}, nil
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}

func (r *PostgresRepository) InsertOHLC(ctx context.Context, ohlc schema.OHLCSchema) error {
	_, err := r.db.Exec(
		queryInsert,
		ohlc.ID,
		ohlc.QueryTime,
		ohlc.Market,
		ohlc.Pair,
		ohlc.Period,
		ohlc.OHLC)
	return err
}

func (r *PostgresRepository) ListOHLCs(ctx context.Context, skip uint64, take uint64) ([]schema.OHLCSchema, error) {
	rows, err := r.db.Query(queryList, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ohlcs := make([]schema.OHLCSchema, 0)
	for rows.Next() {
		ohlc := schema.OHLCSchema{}
		if err = rows.Scan(queryInsert,
			&ohlc.ID,
			&ohlc.QueryTime,
			&ohlc.Market,
			&ohlc.Pair,
			&ohlc.Period,
			&ohlc.OHLC); err == nil {
			ohlcs = append(ohlcs, ohlc)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ohlcs, nil
}
