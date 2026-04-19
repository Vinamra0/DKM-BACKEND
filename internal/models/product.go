package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Slug           string             `bson:"slug" json:"slug"`
	Name           string             `bson:"name" json:"name"`
	Description    string             `bson:"description" json:"description"`
	Image          string             `bson:"image" json:"image"`
	Composition    string             `bson:"composition" json:"composition"`
	DosageForm     string             `bson:"dosageForm" json:"dosageForm"`
	Packing        string             `bson:"packing" json:"packing"`
	Company        string             `bson:"company" json:"company"`
	Category       string             `bson:"category" json:"category"`
	SubCategory    string             `bson:"subCategory" json:"subCategory"`
	PackageType    string             `bson:"packageType" json:"packageType"`
	Tags           []string           `bson:"tags" json:"tags"`
	Generics       []string           `bson:"generics" json:"generics"`
	Specifications map[string]any     `bson:"specifications" json:"specifications"`
	IsActive       bool               `bson:"isActive" json:"isActive"`
	Timestamps     `bson:",inline" json:",inline"`
}
