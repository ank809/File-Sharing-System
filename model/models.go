package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type File struct {
	ID       primitive.ObjectID `bson:"_id"`
	Filename string             `json:"filename"`
	Password string             `json:"password"`
	Url      string             `json:"url"`
}
