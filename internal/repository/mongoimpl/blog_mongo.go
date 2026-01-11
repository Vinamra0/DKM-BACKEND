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

type BlogRepository struct {
	col *mongo.Collection
}

func NewBlogRepository(db *mongo.Database) *BlogRepository {
	return &BlogRepository{col: db.Collection("blogs")}
}

func (r *BlogRepository) FindAll(ctx context.Context) ([]models.Blog, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []models.Blog
	for cur.Next(ctx) {
		var b models.Blog
		if err := cur.Decode(&b); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, cur.Err()
}

func (r *BlogRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Blog, error) {
	var b models.Blog
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&b)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &b, err
}

func (r *BlogRepository) FindBySlug(ctx context.Context, slug string) (*models.Blog, error) {
	var b models.Blog
	err := r.col.FindOne(ctx, bson.M{"slug": slug}).Decode(&b)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &b, err
}

func (r *BlogRepository) Create(ctx context.Context, b *models.Blog) error {
	b.Timestamps.CreatedAt = time.Now()
	b.Timestamps.UpdatedAt = b.Timestamps.CreatedAt
	res, err := r.col.InsertOne(ctx, b)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		b.ID = oid
	}
	return nil
}

func (r *BlogRepository) Update(ctx context.Context, id primitive.ObjectID, b *models.Blog) error {
	b.Timestamps.UpdatedAt = time.Now()
	_, err := r.col.UpdateByID(ctx, id, bson.M{"$set": b})
	return err
}

func (r *BlogRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
