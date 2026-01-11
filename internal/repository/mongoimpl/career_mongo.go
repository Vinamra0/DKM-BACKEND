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

type CareerRepository struct{ col *mongo.Collection }

func NewCareerRepository(db *mongo.Database) *CareerRepository {
	return &CareerRepository{col: db.Collection("careers")}
}

func (r *CareerRepository) FindAll(ctx context.Context) ([]models.Career, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []models.Career
	for cur.Next(ctx) {
		var c models.Career
		if err := cur.Decode(&c); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, cur.Err()
}

func (r *CareerRepository) FindPublic(ctx context.Context) ([]models.Career, error) {
	cur, err := r.col.Find(ctx, bson.M{"active": true})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []models.Career
	for cur.Next(ctx) {
		var c models.Career
		if err := cur.Decode(&c); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, cur.Err()
}

func (r *CareerRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Career, error) {
	var c models.Career
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &c, err
}

func (r *CareerRepository) Create(ctx context.Context, c *models.Career) error {
	if c.Active == false {
		c.Active = true
	}
	c.Timestamps.CreatedAt = time.Now()
	c.Timestamps.UpdatedAt = c.Timestamps.CreatedAt
	res, err := r.col.InsertOne(ctx, c)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		c.ID = oid
	}
	return nil
}

func (r *CareerRepository) Update(ctx context.Context, id primitive.ObjectID, c *models.Career) error {
	c.Timestamps.UpdatedAt = time.Now()
	_, err := r.col.UpdateByID(ctx, id, bson.M{"$set": c})
	return err
}

func (r *CareerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
