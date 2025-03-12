package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `json:"username" binding:"required"`
	Password string             `json:"password" binding:"required"`
}
