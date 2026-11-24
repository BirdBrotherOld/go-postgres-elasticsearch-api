package repository

import (
	"belajar-go-pg-es/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, p *models.Product) error {
	query := `INSERT INTO products (name, description, price) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(ctx, query, p.Name, p.Description, p.Price).Scan(&p.ID)
}

func (r *PostgresRepo) Update(ctx context.Context, p models.Product) (int64, error) {
	query := `UPDATE products SET name=$1, description=$2, price=$3 WHERE id=$4`
	cmdTag, err := r.db.Exec(ctx, query, p.Name, p.Description, p.Price, p.ID)
	if err != nil {
		return 0, err
	}
	return cmdTag.RowsAffected(), nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id int) (int64, error) {
	cmdTag, err := r.db.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return 0, err
	}
	return cmdTag.RowsAffected(), nil
}