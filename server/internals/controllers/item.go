package controllers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jrjaro18/tryingDC/internals/database"
	"github.com/jrjaro18/tryingDC/internals/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllItems(c *fiber.Ctx) error {
	var items []models.Item
	cursor, err := database.Item.Find(context.Background(), bson.M{})
	//store all the items in the items array
	for cursor.Next(context.Background()) {
		var item models.Item
		cursor.Decode(&item)
		items = append(items, item)
	}

	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.JSON(items)
}