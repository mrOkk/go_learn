package postgres

import (
	"App/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MerchantRepository struct {
	pool *pgxpool.Pool
}

func NewMerchantRepository(pool *pgxpool.Pool) *MerchantRepository {
	return &MerchantRepository{
		pool: pool,
	}
}

func (r *MerchantRepository) GetById(ctx context.Context, id int64) (domain.Merchant, error) {
	var merchant domain.Merchant

	row := r.pool.QueryRow(ctx,
		`SELECT id, name, is_active, balance FROM merchants WHERE id = $1`,
		id)

	err := row.Scan(&merchant.Id, &merchant.Name, &merchant.IsActive, &merchant.Balance)
	if err != nil {
		return domain.Merchant{}, fmt.Errorf("get merchant: %w", err)
	}

	return merchant, nil
}

func (r *MerchantRepository) UpdateBalance(ctx context.Context, id int64, newBalance float64) error {
	result, err := executorFrom(ctx, r.pool).Exec(ctx,
		`UPDATE merchants SET balance = $1 WHERE id = $2`,
		newBalance, id)

	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("merchant not found: id=%d", id)
	}

	return nil
}

