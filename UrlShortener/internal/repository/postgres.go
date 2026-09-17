package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *Postgres {
	return &Postgres{
		pool: pool,
	}
}

type ConflictError struct {
	Code string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict code %s", e.Code)
}

func (p *Postgres) Save(ctx context.Context, code string, url string) error {
	_, err := p.pool.Exec(ctx,`INSERT INTO urls (short_code, original_url) VALUES ($1, $2)`, code, url)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, pgErr) && pgErr.Code == "23505" {
			return &ConflictError{Code: code}
		}
		return err
	}
	return nil
}

func (p *Postgres) Get(ctx context.Context, code string) (string, error) {
	var url string
	err := p.pool.QueryRow(
		ctx,
		`SELECT original_url FROM urls WHERE short_code = $1`,
		code,
	).Scan(&url)
	return url, err
}