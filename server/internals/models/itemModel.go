package models
import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Price    int                `json:"price" bson:"price"`
	SellerID primitive.ObjectID `json:"sellerID" bson:"sellerID"`
}