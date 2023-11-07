package db

import (
	"context"
	"github.com/digkill/giggle/schema"
)

type Repository interface {
	Close()
	InsertGiggle(ctx context.Context, giggle schema.Giggle) error
	ListGiggles(ctx context.Context, skip uint64, take uint64) ([]schema.Giggle, error)
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertGiggle(ctx context.Context, giggle schema.Giggle) error {
	return impl.InsertGiggle(ctx, giggle)
}

func ListGiggles(ctx context.Context, skip uint64, take uint64) ([]schema.Giggle, error) {
	return impl.ListGiggles(ctx, skip, take)
}
