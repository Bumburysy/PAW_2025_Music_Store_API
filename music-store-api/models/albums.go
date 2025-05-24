package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Album struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title  string             `bson:"title" json:"title"`
	Artist string             `bson:"artist" json:"artist"`
	Genre  string             `bson:"genre" json:"genre"`
	Price  float64            `bson:"price" json:"price"`
	Stock  int                `bson:"stock" json:"stock"`
}
