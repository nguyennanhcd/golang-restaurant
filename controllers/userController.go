package controller

import (
	"context"
	"golang-restaurant-management/database"
	helper "golang-restaurant-management/helper"
	"golang-restaurant-management/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

type ErrorResponse struct {
	Error string `json:"error"`
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10000)

		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		if s := c.Query("startIndex"); s != "" {
			startIndex, err = strconv.Atoi(s)
			if err != nil {
				startIndex = (page - 1) * recordPerPage
			}
		}

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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		// convert the JSON data coming from postman to something that golang understands

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validate the data based on the struct
		var validate = validator.New()
		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// check if the email already exists in the database

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking email"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}

		// hash the password

		password := HashPassword(*user.Password)
		user.Password = &password

		// you'll also check if the phone no. has already been used by another account

		countPhone, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking phone number"})
			return
		}

		if countPhone > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists"})
			return
		}

		// create some extra details - createdAt, updatedAt, etc, ID

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// generate token and refresh token (gen all tokens from helper function )

		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		// if all ok, then you insert this new user into the user collection

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := "User item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// return status OK and send the result back

		c.JSON(http.StatusOK, gin.H{
			"user_id":       user.User_id,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		// convert the login data from postman which is in JSON format to golang readable format

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// find a user with the email and see if it exists

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		// then verify the password

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		// if all goes well, generate the tokens

		token, refreshToken, generateErr := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

		if generateErr != nil {
			log.Printf("Error generating tokens: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		// update tokens - token and refresh token

		if err := helper.UpdateAllTokens(token, refreshToken, foundUser.User_id); err != nil {
			log.Printf("Error updating tokens: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}

		// return status okn
		c.JSON(http.StatusOK, gin.H{
			"user_id":       foundUser.User_id,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var request struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		claims, msg := helper.ValidateToken(request.RefreshToken)
		if msg != "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: msg})
			return
		}

		if claims.TokenType != "refresh" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid token type, refresh token required"})
			return
		}

		var user struct {
			RefreshToken string `bson:"refresh_token"`
		}

		err := userCollection.FindOne(ctx, bson.M{"user_id": claims.Uid}).Decode(&user)
		if err != nil || user.RefreshToken != request.RefreshToken {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid or revoked refresh token"})
			return
		}

		token, refreshToken, err := helper.GenerateAllTokens(claims.Email, claims.First_name, claims.Last_name, claims.Uid)
		if err != nil {
			log.Printf("Error generating tokens: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate tokens"})
			return
		}

		if err := helper.UpdateAllTokens(token, refreshToken, claims.Uid); err != nil {
			log.Printf("Error updating tokens: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update tokens"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":       claims.Uid,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providePassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providePassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "passwords do not match"
		check = false
	}

	return check, msg
}
