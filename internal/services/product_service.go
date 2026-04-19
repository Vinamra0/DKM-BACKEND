package services

import (
	"context"
	"errors"
	"strings"

	"dkmbackend/internal/models"
	"dkmbackend/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct{ repo repository.ProductRepository }

func NewProductService(r repository.ProductRepository) *ProductService {
	return &ProductService{repo: r}
}

// validateProductFields checks that required specification fields are not empty
func validateProductFields(p *models.Product) error {
	if strings.TrimSpace(p.Composition) == "" {
		return errors.New("composition is required and cannot be empty")
	}
	if strings.TrimSpace(p.DosageForm) == "" {
		return errors.New("dosageForm is required and cannot be empty")
	}
	if strings.TrimSpace(p.Packing) == "" {
		return errors.New("packing is required and cannot be empty")
	}
	if strings.TrimSpace(p.Description) == "" {
		return errors.New("description is required and cannot be empty")
	}
	return nil
}

func (s *ProductService) List(ctx context.Context) ([]models.Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *ProductService) Get(ctx context.Context, id string) (*models.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	return s.repo.FindByID(ctx, oid)
}

func (s *ProductService) Create(ctx context.Context, p *models.Product) error {
	if err := validateProductFields(p); err != nil {
		return err
	}
	return s.repo.Create(ctx, p)
}

func (s *ProductService) Update(ctx context.Context, id string, p *models.Product) error {
	if err := validateProductFields(p); err != nil {
		return err
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	return s.repo.Update(ctx, oid, p)
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	return s.repo.Delete(ctx, oid)
}
