package repository

import (
	"context"

	"github.com/Meystergod/online-store/internal/domain/model"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *model.Product) (string, error)
	GetProduct(ctx context.Context, uuid string) (*model.Product, error)
	GetProductByTitle(ctx context.Context, title string) (*model.Product, error)
	GetAllProducts(ctx context.Context) (*[]model.Product, error)
	UpdateProduct(ctx context.Context, product *model.Product) error
	DeleteProduct(ctx context.Context, uuid string) error
}

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *model.Category) (string, error)
	GetCategory(ctx context.Context, uuid string) (*model.Category, error)
	GetCategoryByTitle(ctx context.Context, title string) (*model.Category, error)
	GetAllCategories(ctx context.Context) (*[]model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category) error
	DeleteCategory(ctx context.Context, uuid string) error
}

type DiscountRepository interface {
	CreateDiscount(ctx context.Context, discount *model.Discount) (string, error)
	GetDiscount(ctx context.Context, uuid string) (*model.Discount, error)
	GetDiscountByTitle(ctx context.Context, title string) (*model.Discount, error)
	GetAllDiscounts(ctx context.Context) (*[]model.Discount, error)
	UpdateDiscount(ctx context.Context, discount *model.Discount) error
	DeleteDiscount(ctx context.Context, uuid string) error
}

type TagRepository interface {
	CreateTag(ctx context.Context, tag *model.Tag) (string, error)
	GetTag(ctx context.Context, uuid string) (*model.Tag, error)
	GetTagByTitle(ctx context.Context, title string) (*model.Tag, error)
	GetAllTags(ctx context.Context) (*[]model.Tag, error)
	UpdateTag(ctx context.Context, tag *model.Tag) error
	DeleteTag(ctx context.Context, uuid string) error
}
