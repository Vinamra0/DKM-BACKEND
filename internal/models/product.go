package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"`
	Image       string             `bson:"image" json:"image"`
	Price       float64            `bson:"price" json:"price"`
	Timestamps  `bson:",inline" json:",inline"`
}
