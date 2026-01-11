package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Career struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Location    string             `bson:"location" json:"location"`
	Type        string             `bson:"type" json:"type"`
	Description string             `bson:"description" json:"description"`
	Active      bool               `bson:"active" json:"active"`
	Timestamps  `bson:",inline" json:",inline"`
}
