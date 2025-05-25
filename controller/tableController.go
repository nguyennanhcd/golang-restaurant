package controller

import (
	"context"
	"golang-restaurant-management/database"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "tables")

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := tableCollection.Find(ctx, bson.D{})

		defer cancel()

		if err != nil {
			c.JSON(500, gin.H{"error": "Error occurred while fetching tables"})
			return
		}

		var allTables []bson.M
		if err = result.All(ctx, &allTables); err != nil {
			log.Fatal(err)
			return
		}

		if len(allTables) == 0 {
			c.JSON(404, gin.H{"message": "No orders found"})
			return
		}

		c.JSON(200, allTables)
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
