package services

import (
	"context"
	"testing"

	"dkmbackend/internal/models"
	"dkmbackend/internal/repository"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type fakeProductRepo struct{ items []models.Product }

var _ repository.ProductRepository = (*fakeProductRepo)(nil)

func (f *fakeProductRepo) FindAll(ctx context.Context) ([]models.Product, error) { return f.items, nil }
func (f *fakeProductRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	for i := range f.items {
		if f.items[i].ID == id {
			return &f.items[i], nil
		}
	}
	return nil, nil
}
func (f *fakeProductRepo) Create(ctx context.Context, p *models.Product) error {
	f.items = append(f.items, *p)
	return nil
}
func (f *fakeProductRepo) Update(ctx context.Context, id primitive.ObjectID, p *models.Product) error {
	return nil
}
func (f *fakeProductRepo) Delete(ctx context.Context, id primitive.ObjectID) error { return nil }

func TestProductService_List(t *testing.T) {
	repo := &fakeProductRepo{items: []models.Product{{Name: "A"}, {Name: "B"}}}
	svc := &ProductService{repo: repo}
	got, err := svc.List(context.Background())
	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "A", got[0].Name)
}
