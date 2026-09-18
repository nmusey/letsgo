package database

import (
	"context"

	"github.com/uptrace/bun"
)

type Repository[T any] struct {
	db *bun.DB
}

func NewRepository[T any](db *bun.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) Get(ctx context.Context, id any) (*T, error) {
	model := new(T)

	err := r.db.NewSelect().
		Model(model).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *Repository[T]) List(ctx context.Context) ([]T, error) {
	var models []T

	err := r.db.NewSelect().Model(&models).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (r *Repository[T]) Create(ctx context.Context, model *T) error {
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	return err
}

func (r *Repository[T]) Update(ctx context.Context, model *T) error {
	_, err := r.db.NewUpdate().Model(model).WherePK().Exec(ctx)
	return err
}

func (r *Repository[T]) Delete(ctx context.Context, model *T) error {
	_, err := r.db.NewDelete().Model(model).WherePK().Exec(ctx)
	return err
}
