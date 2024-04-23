package controllers

import (
	"context"
	"fmt"
	"net/rpc"
	"time"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/jrjaro18/tryingDC/internals/database"
	"github.com/jrjaro18/tryingDC/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func GetAllUsers(c *fiber.Ctx) error {
	//get all users without their passwords
	var users []models.User
	cursor, err := database.User.Find(context.Background(), bson.M{})
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error getting users",
		})
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user models.User
		cursor.Decode(&user)
		user.Password = ""
		users = append(users, user)
	}
	return c.JSON(users)
}

func CreateUser(c *fiber.Ctx) error {
	// send the incoming request to rpc server to create a user
	client, err := rpc.DialHTTP("tcp", "172.16.40.205:1234")
	if err != nil {
		fmt.Println("Error in Dialing")
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error in Dialing",
		})
	}
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error parsing user",
		})
	}
	var reply string

	data := models.LamportRequest{User: *user,  LamportTime: models.MainServerLamportTime}
	fmt.Println(data)
	err = client.Call("API.CreateUser", data, &reply)

	models.MainServerLamportTime++
	fmt.Println("Main Server Updated Lamport Time: ", models.MainServerLamportTime)
	if err != nil {
		fmt.Println("Error in Calling")
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error in Calling",
		})
	}
	if reply != "User created" {
		return c.Status(400).JSON(fiber.Map{
			"message": reply,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": reply,
	})
}

func UserLogin(c *fiber.Ctx) error {
	// login a user
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Error parsing user",
		})
	}
	// if username or email or password is empty then return error
	if user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email and Password are required",
		})
	}
	// find the user with the email
	var result models.User
	err := database.User.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&result)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User with this email does not exist",
		})
	}
	// check the password
	if result.Password != user.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}
	// remove the password from the response
	result.Password = ""
	return c.Status(fiber.StatusAccepted).JSON(result)
}

func AddToCart(c *fiber.Ctx) error {
	// in the body of the data youll have the email of the user and the id of the item to be added to the cart
	data := new(struct {
		Email  string `json:"email"`
		ItemID string `json:"itemID"`
	})
	if err := c.BodyParser(data); err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Error parsing data",
		})
	}
	// find the user with the email
	var result models.User
	err := database.User.FindOne(context.Background(), bson.M{"email": data.Email}).Decode(&result)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User with this email does not exist",
		})
	}
	// add the item to the cart
	//if cart is nil then create a new array
	if result.Cart == nil {
		result.Cart = []primitive.ObjectID{}
	}
	// convert the string to object id
	itemID, err := primitive.ObjectIDFromHex(data.ItemID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid itemID",
		})
	}
	// add the item to the cart
	result.Cart = append(result.Cart, itemID)
	// update the user
	_, err = database.User.UpdateOne(context.Background(), bson.M{"email": data.Email}, bson.M{"$set": bson.M{"cart": result.Cart}})
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error adding item to cart",
		})
	}
	// remove the password from the response
	result.Password = ""
	return c.Status(fiber.StatusAccepted).JSON(result, "Item added to cart")
}

func RemoveFromCart(c *fiber.Ctx) error {
	// in the body of the data youll have the email of the user and the id of the item to be removed from the cart
	data := new(struct {
		Email  string `json:"email"`
		ItemID string `json:"itemID"`
	})
	if err := c.BodyParser(data); err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Error parsing data",
		})
	}
	// find the user with the email
	var result models.User
	err := database.User.FindOne(context.Background(), bson.M{"email": data.Email}).Decode(&result)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User with this email does not exist",
		})
	}
	// remove the item from the cart
	//if cart is nil then return error
	if result.Cart == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cart is empty",
		})
	}
	// convert the string to object id
	itemID, err := primitive.ObjectIDFromHex(data.ItemID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid itemID",
		})
	}
	// remove the item from the cart
	for i, v := range result.Cart {
		if v == itemID {
			result.Cart = append(result.Cart[:i], result.Cart[i+1:]...)
			break
		}
	}
	// update the user
	_, err = database.User.UpdateOne(context.Background(), bson.M{"email": data.Email}, bson.M{"$set": bson.M{"cart": result.Cart}})
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error removing item from cart",
		})
	}
	// remove the password from the response
	result.Password = ""
	return c.Status(fiber.StatusAccepted).JSON(result)
}

func GetCart(c *fiber.Ctx) error {
	// in the body of the data youll have the email of the user
	data := new(struct {
		Email string `json:"email"`
	})
	if err := c.BodyParser(data); err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Error parsing data",
		})
	}
	// find the user with the email
	var result models.User
	err := database.User.FindOne(context.Background(), bson.M{"email": data.Email}).Decode(&result)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User with this email does not exist",
		})
	}
	// remove the password from the response
	result.Password = ""
	
	var items = make([]models.Item, len(result.Cart))
	for i, v := range result.Cart {
		var item models.Item
		err := database.Item.FindOne(context.Background(), bson.M{"_id": v}).Decode(&item)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Item not found",
			})
		}
		items[i] = item
	}
	return c.Status(fiber.StatusAccepted).JSON(items)
}

func Buy(c *fiber.Ctx) error {
	// in the body of the data youll have the email of the user

	data := new(struct {
		Email string `json:"email"`
	})
	if err := c.BodyParser(data); err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Error parsing data",
		})
	}
	// find the user with the email
	var result models.User
	err := database.User.FindOne(context.Background(), bson.M{"email": data.Email}).Decode(&result)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User with this email does not exist",
		})
	}

	if len(result.Cart)==0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The Cart is Empty!🛒",
		})
	}

	// add the items in the cart to the previously bought
	//if previously bought is nil then create a new array
	if result.PreviouslyBought == nil {
		result.PreviouslyBought = []primitive.ObjectID{}
	}
	// add the items to the previously bought
	result.PreviouslyBought = append(result.PreviouslyBought, result.Cart...)
	// update the user
	_, err = database.User.UpdateOne(context.Background(), bson.M{"email": data.Email}, bson.M{"$set": bson.M{"previouslyBought": result.PreviouslyBought}})
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error buying items",
		})
	}

	x := result.ID.String() + " bought "
	for _, temp := range(result.Cart) {
		x = x + temp.String() + ", "
	}
	x = x + " at " + time.Now().String()

	go AddToLogs(x)
	
	// remove the items from the cart
	result.Cart = []primitive.ObjectID{}
	// update the user
	_, err = database.User.UpdateOne(context.Background(), bson.M{"email": data.Email}, bson.M{"$set": bson.M{"cart": result.Cart}})
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error buying items",
		})
	}

	// remove the password from the response
	result.Password = ""	
	return c.Status(fiber.StatusAccepted).JSON(result)
}

func PreviouslyBought(c *fiber.Ctx) error {
	// in the body of the data youll have the email of the user
	data := new(struct {
		Email string `json:"email"`
	})
	if err := c.BodyParser(data); err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Error parsing data",
		})
	}
	// find the user with the email
	var result models.User
	err := database.User.FindOne(context.Background(), bson.M{"email": data.Email}).Decode(&result)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User with this email does not exist",
		})
	}
	//fetch all the items from the previously bought array
	var items = make([]models.Item, len(result.PreviouslyBought))

	for i, v := range result.PreviouslyBought {
		var item models.Item
		err := database.Item.FindOne(context.Background(), bson.M{"_id": v}).Decode(&item)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Item not found",
			})
		}
		items[i] = item
	}
	// remove the password from the response
	result.Password = ""
	return c.Status(fiber.StatusAccepted).JSON(items)
}

func AddToLogs(s string) error {
	database.Mutex.Lock()
	time.Sleep(10 * time.Second)

	fileToWrite, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return err
    }
    defer fileToWrite.Close()

    textToAppend := "\n"+s
    _, err = fileToWrite.WriteString(textToAppend)
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return err
    }
    fmt.Println("Append successful!")
	database.Mutex.Unlock()
	return nil
}

