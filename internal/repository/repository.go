package repository

import (
	"context"

	"dkmbackend/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRepository interface {
	FindAll(ctx context.Context) ([]models.Product, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error)
	Create(ctx context.Context, p *models.Product) error
	Update(ctx context.Context, id primitive.ObjectID, p *models.Product) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type BlogRepository interface {
	FindAll(ctx context.Context) ([]models.Blog, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Blog, error)
	FindBySlug(ctx context.Context, slug string) (*models.Blog, error)
	Create(ctx context.Context, b *models.Blog) error
	Update(ctx context.Context, id primitive.ObjectID, b *models.Blog) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type CareerRepository interface {
	FindAll(ctx context.Context) ([]models.Career, error)
	FindPublic(ctx context.Context) ([]models.Career, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Career, error)
	Create(ctx context.Context, c *models.Career) error
	Update(ctx context.Context, id primitive.ObjectID, c *models.Career) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
