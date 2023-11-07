package db

import (
	"context"
	"database/sql"

	"github.com/digkill/giggle/schema"
	_ "github.com/lib/pq"
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

func (r *PostgresRepository) InsertGiggle(ctx context.Context, giggle schema.Giggle) error {
	_, err := r.db.Exec("INSERT INTO giggle(id, body, created_at) VALUES($1, $2, $3)", giggle.ID, giggle.Body, giggle.CreatedAt)
	return err
}

func (r *PostgresRepository) ListGiggles(ctx context.Context, skip uint64, take uint64) ([]schema.Giggle, error) {
	rows, err := r.db.Query("SELECT * FROM giggle ORDER BY id DESC OFFSET $1 LIMIT $2", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	giggles := []schema.Giggle{}
	for rows.Next() {
		giggle := schema.Giggle{}
		if err = rows.Scan(&giggle.ID, &giggle.Body, &giggle.CreatedAt); err == nil {
			giggles = append(giggles, giggle)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return giggles, nil
}
