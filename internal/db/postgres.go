package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionConfig struct {
	Host     string
	Port     string
	DbName   string
	Username string
	Password string
}

func (conf *ConnectionConfig) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.DbName)
}

var dbPool *pgxpool.Pool

func Connect(ctx context.Context, config ConnectionConfig) error {
	if dbPool != nil {
		return errors.New("db connections pool is already established")
	}
	dbpool, err := pgxpool.New(ctx, config.String())
	if err != nil {
		return err
	}
	dbPool = dbpool
	return nil
}

func Disconnect() error {
	if dbPool == nil {
		return errors.New("db connection was never established. No need to call disconnect")
	}
	dbPool.Close()
	return nil
}

func GetPool() (*pgxpool.Pool, error) {
	if dbPool == nil {
		return nil, errors.New("pool is not established")
	}
	return dbPool, nil
}
