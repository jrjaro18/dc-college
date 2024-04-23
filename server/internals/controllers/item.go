package controllers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jrjaro18/tryingDC/internals/database"
	"github.com/jrjaro18/tryingDC/internals/models"
	"github.com/jrjaro18/tryingDC/internals/redis"
	"go.mongodb.org/mongo-driver/bson"
	"encoding/json"
)

func GetAllItems(c *fiber.Ctx) error {
	cachedItems, err := redis.RedisClient.Get(context.Background(), "all_items").Result()
    if err == nil {
        var items []models.Item
        json.Unmarshal([]byte(cachedItems), &items)
		fmt.Println("Retrieved items from cache")
        return c.JSON(items)
    }
	var items []models.Item
	cursor, err := database.Item.Find(context.Background(), bson.M{})
	//store all the items in the items array
	count := 0
	for cursor.Next(context.Background()) {
		var item models.Item
		cursor.Decode(&item)
		items = append(items, item)
		count++
	}

	if err != nil {
		fmt.Println(err)
		return err
	}
	if count >= 5 {
		// return only the last 5
		items = items[count-5:]
	}
	itemsJSON, _ := json.Marshal(items)
    redis.RedisClient.Set(context.Background(), "all_items", string(itemsJSON), 0)

	return c.JSON(items)
}