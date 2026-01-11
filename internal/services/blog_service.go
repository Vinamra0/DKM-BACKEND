package services

import (
	"context"
	"errors"

	"dkmbackend/internal/models"
	"dkmbackend/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogService struct{ repo repository.BlogRepository }

func NewBlogService(r repository.BlogRepository) *BlogService { return &BlogService{repo: r} }

func (s *BlogService) List(ctx context.Context) ([]models.Blog, error) { return s.repo.FindAll(ctx) }

func (s *BlogService) GetByParam(ctx context.Context, param string) (*models.Blog, error) {
	if oid, err := primitive.ObjectIDFromHex(param); err == nil {
		return s.repo.FindByID(ctx, oid)
	}
	b, err := s.repo.FindBySlug(ctx, param)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("not found")
	}
	return b, nil
}

func (s *BlogService) Create(ctx context.Context, b *models.Blog) error { return s.repo.Create(ctx, b) }
func (s *BlogService) Update(ctx context.Context, id string, b *models.Blog) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	return s.repo.Update(ctx, oid, b)
}
func (s *BlogService) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	return s.repo.Delete(ctx, oid)
}
