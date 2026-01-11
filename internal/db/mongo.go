package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)
	return mongo.Connect(ctx, opts)
}
