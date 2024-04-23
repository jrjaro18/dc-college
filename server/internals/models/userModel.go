package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
	Cart    []primitive.ObjectID `json:"cart" bson:"cart"`
	PreviouslyBought []primitive.ObjectID `json:"previouslyBought" bson:"previouslyBought"`
}

type LamportRequest struct {
	User User
	Seller Seller
	LamportTime int
}

var MainServerLamportTime = 0
