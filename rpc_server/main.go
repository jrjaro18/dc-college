package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var User1 *mongo.Collection
var Seller1 *mongo.Collection

func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://rohan:rohan@cluster0.piveseb.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}

	// Initialize collections
	Seller1 = client.Database("dc-ecommerce").Collection("sellers")
	User1 = client.Database("dc-ecommerce").Collection("users")

	// Ping MongoDB
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Error pinging MongoDB:", err)
		return
	}

	// Register API methods
	apiRPC := new(API)
	err = rpc.Register(apiRPC)
	if err != nil {
		fmt.Println("Error in Registering API:", err)
		return
	}
	rpc.HandleHTTP()

	// Start RPC servers
	go startRPCServer(":1234")

	// Keep the main goroutine running
	select {}
}

// startRPCServer starts an RPC server listening on the specified port
func startRPCServer(port string) {
	fmt.Println("Starting RPC Server on port", port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error in Listening on port", port, ":", err)
		return
	}
	defer listener.Close()

	err = http.Serve(listener, nil)
	if err != nil {
		fmt.Println("Error serving RPC requests on port", port, ":", err)
	}
}


type API int

var rpcServerLamportTime = 0

func (a *API) CreateSeller(req LamportRequest, reply *string) error {
	// if username or email or password is empty then return error
	seller := req.Seller
	reqTime := req.LamportTime

	fmt.Printf("Request Time: %d, RPC Server Time: %d\n", reqTime, rpcServerLamportTime)

	if reqTime > rpcServerLamportTime {
		rpcServerLamportTime = reqTime
	}
	rpcServerLamportTime++
	fmt.Println("Updated RPC Server Time: ", rpcServerLamportTime)

	if seller.Username == "" || seller.Email == "" || seller.Password == "" {
		fmt.Printf("Username: %s, Email: %s, Password: %s\n", seller.Username, seller.Email, seller.Password)
		*reply = "Username, Email or Password is empty"
		return nil

	}

	// if the person with same email is in the database then return error
	var result User
	err := Seller1.FindOne(context.Background(), bson.M{"email": seller.Email}).Decode(&result)
	if err == nil {
		*reply = "Email already exists"
		return nil
	}

	// otherwise create the user
	_, err = Seller1.InsertOne(context.Background(), seller)
	if err != nil {
		fmt.Println(err)
		*reply = "Error in creating the seller"
	}
	// remove the password from the response
	seller.Password = ""
	*reply = "Seller created"
	return nil
}

func (a *API) CreateUser(req LamportRequest, reply *string) error {

	user := req.User
	reqTime := req.LamportTime

	fmt.Printf("Request Time: %d, RPC Server Time: %d\n", reqTime, rpcServerLamportTime)

	if reqTime > rpcServerLamportTime {
		rpcServerLamportTime = reqTime
	}
	rpcServerLamportTime++
	fmt.Println("Updated RPC Server Time: ", rpcServerLamportTime)

	// if username or email or password is empty then return error
	if user.Username == "" || user.Email == "" || user.Password == "" {
		fmt.Printf("Username: %s, Email: %s, Password: %s\n", user.Username, user.Email, user.Password)
		*reply = "Username, Email or Password is empty"
		return nil
	}

	// if the person with same email is in the database then return error
	var result User
	err := User1.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&result)
	if err == nil {
		*reply = "Email already exists"
		return nil
	}

	// otherwise create the user
	_, err = User1.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println(err)
		*reply = "Error in creating the user"
		return nil
	}
	// remove the password from the response
	user.Password = ""
	*reply = "User created"
	return nil
}

type Item struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Price    int                `json:"price" bson:"price"`
	SellerID primitive.ObjectID `json:"sellerID" bson:"sellerID"`
}

type Seller struct {
	ID       primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string               `json:"username" bson:"username"`
	Password string               `json:"password" bson:"password"`
	Email    string               `json:"email" bson:"email"`
	Items    []primitive.ObjectID `json:"items" bson:"items"`
}

type SellerResponse struct {
	ID       primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string               `json:"username" bson:"username"`
	Email    string               `json:"email" bson:"email"`
	Items    []primitive.ObjectID `json:"items" bson:"items"`
}

type User struct {
	ID               primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Username         string               `json:"username" bson:"username"`
	Password         string               `json:"password" bson:"password"`
	Email            string               `json:"email" bson:"email"`
	Cart             []primitive.ObjectID `json:"cart" bson:"cart"`
	PreviouslyBought []primitive.ObjectID `json:"previouslyBought" bson:"previouslyBought"`
}

type LamportRequest struct {
	User        User
	Seller      Seller
	LamportTime int
}

var MainServerLamportTime = 0
