package rpc

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"net/http"
	"github.com/jrjaro18/tryingDC/internals/database"
	"github.com/jrjaro18/tryingDC/internals/models"
	"go.mongodb.org/mongo-driver/bson"
)

type API int

var rpcServerLamportTime = 0


func Init() {
	fmt.Println("Starting RPC Server")
	var api_rpc = new(API)
	err := rpc.Register(api_rpc)
	if err != nil {
		fmt.Println("Error in Registering API")
		fmt.Println(err.Error())
	}
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error in Listening")
		fmt.Println(err.Error())
	}
	fmt.Println("Listening on port 1234")
	
	http.Serve(listener, nil)
	listener.Close()
	fmt.Println("RPC Server closed")
}

func (a *API) CreateSeller(seller models.Seller, reply *string) error {
	// if username or email or password is empty then return error
	if seller.Username == "" || seller.Email == "" || seller.Password == "" {
		fmt.Printf("Username: %s, Email: %s, Password: %s\n", seller.Username, seller.Email, seller.Password)
		*reply = "Username, Email or Password is empty"
		return nil

	}

	// if the person with same email is in the database then return error
	var result models.User
	err := database.Seller.FindOne(context.Background(), bson.M{"email": seller.Email}).Decode(&result)
	if err == nil {
		*reply = "Email already exists"
		return nil
	}

	// otherwise create the user
	_, err = database.Seller.InsertOne(context.Background(), seller)
	if err != nil {
		fmt.Println(err)
		*reply = "Error in creating the seller"
	}
	// remove the password from the response
	seller.Password = ""
	*reply = "Seller created"
	return nil
}

func (a *API) CreateUser(req models.LamportRequest, reply *string) error {

	user := req.User
	reqTime := req.LamportTime

	fmt.Printf("Request Time: %d, RPC Server Time: %d\n", reqTime, rpcServerLamportTime)

	if(reqTime > rpcServerLamportTime){
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
	var result models.User
	err := database.User.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&result)
	if err == nil {
		*reply = "Email already exists"
		return nil
	}

	// otherwise create the user
	_, err = database.User.InsertOne(context.Background(), user)
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
