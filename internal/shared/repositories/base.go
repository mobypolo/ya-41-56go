package repositories

import (
	"context"

	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*T, error)
	FindByIDWithPreloads(ctx context.Context, id uint, preloads ...string) (*T, error)
	FindAll(ctx context.Context) ([]T, error)
	FindAllPaginated(ctx context.Context, pagination Pagination, preloads ...string) ([]T, bool, error)

	FindByField(ctx context.Context, field string, value any) (*T, error)
	FindByFieldWithPreloads(ctx context.Context, field string, value any, preloads ...string) (*T, error)
	FindManyByField(ctx context.Context, field string, value any) ([]T, error)
}

type GormRepository[T any] struct {
	db *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{db: db}
}

func (r *GormRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *GormRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *GormRepository[T]) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(new(T), "id = ?", id).Error
}

func (r *GormRepository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	return &entity, err
}

func (r *GormRepository[T]) FindByIDWithPreloads(ctx context.Context, id uint, preloads ...string) (*T, error) {
	var entity T
	query := r.db.WithContext(ctx)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.First(&entity, id).Error
	return &entity, err
}

func (r *GormRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

type Pagination struct {
	Page    int
	PerPage int
}

func (r *GormRepository[T]) FindAllPaginated(
	ctx context.Context,
	pagination Pagination,
	preloads ...string,
) ([]T, bool, error) {
	var (
		entities []T
		total    int64
	)

	baseQuery := r.db.WithContext(ctx).Model(new(T))

	for _, preload := range preloads {
		baseQuery = baseQuery.Preload(preload)
	}

	countQuery := baseQuery.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, false, err
	}

	offset := (pagination.Page - 1) * pagination.PerPage
	if err := baseQuery.Offset(offset).Limit(pagination.PerPage).Find(&entities).Error; err != nil {
		return nil, false, err
	}

	hasNext := int64(offset+pagination.PerPage) < total
	return entities, hasNext, nil
}

func (r *GormRepository[T]) FindByField(ctx context.Context, field string, value any) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).Where(field+" = ?", value).First(&entity).Error
	return &entity, err
}

func (r *GormRepository[T]) FindByFieldWithPreloads(ctx context.Context, field string, value any, preloads ...string) (*T, error) {
	var entity T
	query := r.db.WithContext(ctx)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.Where(field+" = ?", value).First(&entity).Error
	return &entity, err
}

func (r *GormRepository[T]) FindManyByField(ctx context.Context, field string, value any) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Where(field+" = ?", value).Find(&entities).Error
	return entities, err
}
