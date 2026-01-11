package services

import (
	"context"
	"errors"

	"dkmbackend/internal/models"
	"dkmbackend/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CareerService struct{ repo repository.CareerRepository }

func NewCareerService(r repository.CareerRepository) *CareerService { return &CareerService{repo: r} }

func (s *CareerService) List(ctx context.Context) ([]models.Career, error) {
	return s.repo.FindAll(ctx)
}
func (s *CareerService) Public(ctx context.Context) ([]models.Career, error) {
	return s.repo.FindPublic(ctx)
}
func (s *CareerService) Get(ctx context.Context, id string) (*models.Career, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	return s.repo.FindByID(ctx, oid)
}
func (s *CareerService) Create(ctx context.Context, c *models.Career) error {
	return s.repo.Create(ctx, c)
}
func (s *CareerService) Update(ctx context.Context, id string, c *models.Career) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	return s.repo.Update(ctx, oid, c)
}
func (s *CareerService) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	return s.repo.Delete(ctx, oid)
}
