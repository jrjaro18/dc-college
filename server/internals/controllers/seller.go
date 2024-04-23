package controllers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/rpc"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jrjaro18/tryingDC/internals/database"
	"github.com/jrjaro18/tryingDC/internals/models"
	"github.com/jrjaro18/tryingDC/internals/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateSeller(c *fiber.Ctx) error {
	// send the request to the rpc server to create a seller
	client, err := rpc.DialHTTP("tcp", "172.16.40.205:1234")
	if err != nil {
		fmt.Println("Error in Dialing")
		fmt.Println(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error in Dialing",
		})
	}

	seller := new(models.Seller)
	if err := c.BodyParser(seller); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error parsing seller",
		})
	}

	var reply string
	// err = client.Call("API.CreateSeller", seller, &reply)

	data := models.LamportRequest{Seller: *seller,  LamportTime: models.MainServerLamportTime}
	fmt.Println(data)
	err = client.Call("API.CreateSeller", data, &reply)

	models.MainServerLamportTime++
	fmt.Println("Main Server Updated Lamport Time: ", models.MainServerLamportTime)
	if err != nil {
		fmt.Println("Error in Calling")
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error in Calling",
		})
	}

	if reply != "Seller created" {
		return c.Status(400).JSON(fiber.Map{
			"message": reply,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": reply,
	})
}

// login seller
func LoginSeller(c *fiber.Ctx) error {
	//login a user
	seller := new(models.Seller)
	if err := c.BodyParser(seller); err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Error parsing seller",
		})
	}
	// if email or password is empty then return error
	if seller.Email == "" || seller.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Email and Password are required",
		})
	}
	// if the person with same email is in the database then return error
	var result models.Seller
	err := database.Seller.FindOne(context.Background(), bson.M{"email": seller.Email}).Decode(&result)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User with this email does not exist",
		})
	}
	// if the password is wrong then return error
	if result.Password != seller.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Password is wrong",
		})
	}
	// remove the password from the response
	result.Password = ""
	return c.Status(fiber.StatusOK).JSON(result)
}

// add item for seller
func AddItem(c *fiber.Ctx) error {
	//add an item to the seller
	item := new(models.Item)
	if err := c.BodyParser(item); err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Error parsing item",
		})
	}
	// if name or price is empty then return error
	if item.Name == "" || item.Price == 0 || item.SellerID == primitive.NilObjectID {
		return c.Status(400).JSON(fiber.Map{
			"message": "Name and Price are required",
		})
	}
	// if the seller with same id is in the database then return error
	var result models.Seller
	err := database.Seller.FindOne(context.Background(), bson.M{"_id": item.SellerID}).Decode(&result)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Seller with this id does not exist",
		})
	}
	// otherwise create the item
	_, err = database.Item.InsertOne(context.Background(), item)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating item",
		})
	}

	// get the item from the database
	err = database.Item.FindOne(context.Background(), bson.M{"name": item.Name, "sellerID": item.SellerID}).Decode(&item)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting item",
		})
	}
	// add the item to the seller
	// if the seller items array is empty then create a new array
	if result.Items == nil {
		result.Items = []primitive.ObjectID{}
	}
	// append the item id to the seller items array
	result.Items = append(result.Items, item.ID)
	// update the seller in the database
	_, err = database.Seller.UpdateOne(context.Background(), bson.M{"_id": item.SellerID}, bson.M{"$set": bson.M{"items": result.Items}})
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating seller",
		})
	}
	// redis
	itemJSON, err := json.Marshal(item)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error marshaling item",
		})
	}
	key := fmt.Sprintf("item:%s", item.ID.Hex())
	err = redis.RedisClient.Set(context.Background(), key, itemJSON, 24*time.Hour).Err() // Cache item for 24 hours
	if err != nil {
		fmt.Println(err)
		// Handle cache error (optional)
	}

	// Get the current value of the "all_items" cache key
	allItemsJSON, err := redis.RedisClient.Get(context.Background(), "all_items").Bytes()
	var allItems []models.Item
	if err == nil {
		// Unmarshal the current value of "all_items" into a slice of items
		err = json.Unmarshal(allItemsJSON, &allItems)
		if err != nil {
			fmt.Println(err)
			// Handle unmarshal error (optional)
		}
	}
	// Append the new item to the slice of items
	allItems = append(allItems, *item)
	if len(allItems) > 5 {
		// Return only the last 5 items
		allItems = allItems[len(allItems)-5:]
	}
	// Marshal the updated slice of items to JSON
	allItemsJSON, err = json.Marshal(allItems)
	if err != nil {
		fmt.Println(err)
		// Handle marshal error (optional)
	}
	// Update the value of the "all_items" cache key
	err = redis.RedisClient.Set(context.Background(), "all_items", allItemsJSON, 24*time.Hour).Err() // Cache item for 24 hours
	if err != nil {
		fmt.Println(err)
		// Handle cache error (optional)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

//func get logs

func GetLogs(c *fiber.Ctx) error {
	database.Mutex.Lock()
	// Reading from a text file
	fileToRead, err := os.Open("logs.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Error parsing item",
		})
	}
	defer fileToRead.Close()

	// Reading line by line
	scanner := bufio.NewScanner(fileToRead)
	x := make([]string, 0)

	for scanner.Scan() {
		y := scanner.Text()
		fmt.Println(y)
		x = append(x, y)
	}

	fmt.Println(x)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	database.Mutex.Unlock()

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": x,
	})
}
