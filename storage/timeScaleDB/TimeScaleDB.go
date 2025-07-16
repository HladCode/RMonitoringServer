package timescaledb

import (
	"context"

	"github.com/HladCode/RMonitoringServer/internal/lib/e"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool  *pgxpool.Pool
	cntxt context.Context
}

func NewDatabase(connStr string) (*Database, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, e.Wrap("can not open DB", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, e.Wrap("can not connect to DB", err)
	}

	return &Database{pool: pool, cntxt: ctx}, nil
}

func (db *Database) InitFromFile(path string) error {
	// TODO просто сделать проверку на наличие таблиц,
	// а то некоторые запросы плохо исполняются по много раз

	// content, err := os.ReadFile(path)
	// if err != nil {
	// 	return e.Wrap("Can not read file", err)
	// }

	// _, err = db.pool.Exec(db.cntxt, string(content))
	// if err != nil {
	// 	return e.Wrap("Exec schema error", err)
	// }

	// хз все можно сделать в бобре

	return nil
}
