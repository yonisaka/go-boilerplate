package datastore

import (
	"context"

	"github.com/yonisaka/go-boilerplate/internal/entities/repository"
	"golang.org/x/sync/errgroup"
)

// healthRepo is a health repository.
type healthRepo struct {
	*BaseRepo
}

// NewHealthRepo returns a health repository.
func NewHealthRepo(base *BaseRepo) repository.HealthRepo {
	return &healthRepo{
		BaseRepo: base,
	}
}

// GetLiveness returns liveness.
func (r *healthRepo) GetLiveness(ctx context.Context) error {
	var err error

	group, _ := errgroup.WithContext(ctx)

	group.Go(func() error {
		return r.dbSlave.Ping(ctx)
	})

	group.Go(func() error {
		_, err = r.dbSlave.Query(ctx, "SELECT 1")
		return err
	})

	group.Go(func() error {
		return r.dbMaster.Ping(ctx)
	})

	group.Go(func() error {
		_, err = r.dbMaster.Query(ctx, "SELECT 1")
		return err
	})

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}
