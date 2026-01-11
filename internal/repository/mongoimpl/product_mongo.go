package mongoimpl

import (
	"context"
	"errors"
	"time"

	"dkmbackend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct{ col *mongo.Collection }

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{col: db.Collection("products")}
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]models.Product, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []models.Product
	for cur.Next(ctx) {
		var p models.Product
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, cur.Err()
}

func (r *ProductRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	var p models.Product
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &p, err
}

func (r *ProductRepository) Create(ctx context.Context, p *models.Product) error {
	p.Timestamps.CreatedAt = time.Now()
	p.Timestamps.UpdatedAt = p.Timestamps.CreatedAt
	res, err := r.col.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		p.ID = oid
	}
	return nil
}

func (r *ProductRepository) Update(ctx context.Context, id primitive.ObjectID, p *models.Product) error {
	p.Timestamps.UpdatedAt = time.Now()
	_, err := r.col.UpdateByID(ctx, id, bson.M{"$set": p})
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
