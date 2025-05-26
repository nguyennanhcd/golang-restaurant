package controller

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10000)

		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {
			recordPerPage = 10 // default value
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1 // default value
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key: "$match", Value: bson.D{}}}

		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$user_items", startIndex, recordPerPage}}}},
			}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			projectStage,
		})

		if err != nil {
			c.JSON(500, gin.H{"error": "Error occurred while fetching users"})
			return
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10000)

		defer cancel()

		userId := c.Param("user_id")

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching user"})
			return
		}

		c.JSON(http.StatusOK, user)

	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// convert the JSON data coming from postman to something that golang understands

		// validate the data based on the struct

		// hash the password

		// you'll also check if the phone no. has already been used by another account

		// create some extra details - createdAt, updatedAt, etc, ID

		// generate token and refresh token (gen all tokens from helper function )

		// if all ok, then you insert this new user into the user collection

		// return status OK and send the result back
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// convert the login data from postman which is in JSON format to golang readable format

		// find a user with the email and see if it exists

		// then verify the password

		// if all goes well, generate the tokens

		// update tokens - token and refresh token

		// return status okn
	}
}

func HashPassword(password string) string {
	return password
}

func VerifyPassword(userPassword string, providePassword string) (bool, string) {
	return true, ""
}
