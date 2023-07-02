package datastore

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yonisaka/go-boilerplate/config"
	"github.com/yonisaka/go-boilerplate/pkg/di"
)

var (
	poolMasterOnce sync.Once
	poolSlaveOnce  sync.Once
	poolMaster     *pgxpool.Pool
	poolSlave      *pgxpool.Pool
)

type wrapPool struct {
	pool *pgxpool.Pool
}

func (w *wrapPool) Close() error {
	w.pool.Close()
	return nil
}

// NewBaseRepo returns a base repository.
func NewBaseRepo(dbMaster, dbSlave *pgxpool.Pool) *BaseRepo {
	return &BaseRepo{dbMaster: dbMaster, dbSlave: dbSlave}
}

// BaseRepo is a base repository.
type BaseRepo struct {
	dbMaster *pgxpool.Pool
	dbSlave  *pgxpool.Pool
}

func getConnString(cfg *config.DB) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host+":"+cfg.Port,
		cfg.DB,
	)
}

// GetDatabaseMaster returns postgresql Pool for Master.
func GetDatabaseMaster(cfg *config.DB) *pgxpool.Pool {
	poolMasterOnce.Do(func() {
		ctx := context.Background()

		var err error

		connString := getConnString(cfg)

		// Use default config.
		poolMaster, err = pgxpool.New(ctx, connString)
		if err != nil {
			log.Fatalf("failed to connect to timescaleDB pool: %v", err)
		}

		err = poolMaster.Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping database: %v", err)
		}

		var c io.Closer = &wrapPool{
			pool: poolMaster,
		}

		di.RegisterCloser("TimescaleDB Master Connection", c)
	})

	return poolMaster
}

// GetDatabaseSlave returns postgresql Pool for Slave.
func GetDatabaseSlave(cfg *config.DB) *pgxpool.Pool {
	poolSlaveOnce.Do(func() {
		ctx := context.Background()

		var err error

		connString := getConnString(cfg)

		// Use default config.
		poolSlave, err = pgxpool.New(ctx, connString)
		if err != nil {
			log.Fatalf("failed to connect to timescaleDB pool: %v", err)
		}

		err = poolSlave.Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping database: %v", err)
		}

		var c io.Closer = &wrapPool{
			pool: poolSlave,
		}

		di.RegisterCloser("TimescaleDB Slave Connection", c)
	})

	return poolSlave
}
