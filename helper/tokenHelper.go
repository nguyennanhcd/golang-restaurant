package helper

import (
	"context"
	"fmt"
	"golang-restaurant-management/database"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	TokenType  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var SECRET_KEY string

func init() {
	SECRET_KEY = os.Getenv("SECRET_KEY")
	if SECRET_KEY == "" {
		log.Fatal("SECRET_KEY environment variable not set")
	}
}

func GenerateAllTokens(email, firstName, lastName, uid string) (signedToken string, refreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		TokenType:  "access",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		TokenType:  "refresh",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24 * time.Duration(30)).Unix(),
		},
	}

	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return signedToken, refreshToken, nil
}

func UpdateAllTokens(signedToken, signedRefreshToken, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: time.Now()})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{Upsert: &upsert}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{{Key: "$set", Value: updateObj}},
		&opt,
	)
	return err
}

func ValidateToken(signedToken string) (*SignedDetails, string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		return nil, fmt.Sprintf("Token validation failed: %v", err)
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok || !token.Valid {
		return nil, "Invalid token"
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, "Token is expired"
	}

	return claims, ""
}
