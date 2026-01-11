package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Blog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title      string             `bson:"title" json:"title"`
	Excerpt    string             `bson:"excerpt" json:"excerpt"`
	Content    string             `bson:"content" json:"content"`
	Author     string             `bson:"author" json:"author"`
	Date       string             `bson:"date" json:"date"`
	Image      string             `bson:"image" json:"image"`
	Category   string             `bson:"category" json:"category"`
	Slug       string             `bson:"slug" json:"slug"`
	Timestamps `bson:",inline" json:",inline"`
}
