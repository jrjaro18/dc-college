package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Seller struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
	Items    []primitive.ObjectID `json:"items" bson:"items"`
}

type SellerResponse	struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Items    []primitive.ObjectID `json:"items" bson:"items"`
}
