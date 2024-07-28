package db

import (
	"context"

	"github.com/digkill/giggle/schema"
)

// Repository interface
type Repository interface {
	Close()
	InsertGiggle(ctx context.Context, giggle schema.Giggle) error
	ListGiggles(ctx context.Context, skip uint64, take uint64) ([]schema.Giggle, error)
}

var repo Repository

func SetRepository(repository Repository) {
	repo = repository
}

func Close() {
	repo.Close()
}

func InsertGiggle(ctx context.Context, giggle schema.Giggle) error {
	return repo.InsertGiggle(ctx, giggle)
}

func ListGiggles(ctx context.Context, skip uint64, take uint64) ([]schema.Giggle, error) {
	return repo.ListGiggles(ctx, skip, take)
}
