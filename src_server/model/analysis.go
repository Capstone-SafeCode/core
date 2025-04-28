package model

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Analysis struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Timestamp primitive.DateTime `bson:"timestamp" json:"timestamp"`
	Results   []gin.H            `bson:"results" json:"results"`
}
