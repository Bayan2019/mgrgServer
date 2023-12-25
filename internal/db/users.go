package db

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	FirstName      string             `json:"first_name"`
	LastName       string             `json:"last_name"`
	Email          string             `json:"email"`
	Phone          string             `json:"phone"`
	Role           string             `json:"role"`
	HashedPassword string             `json:"hashed_password"`
	Verified       bool               `json:"verified"`
}

var userCollection *mongo.Collection = OpenCollection(Client, "users")
var ErrAlreadyExists = errors.New("already exists")

func CreateUser(firstName string, lastName string, email string, phone string, role string, password string) (User, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	notVerified := false

	user := User{
		ID:             primitive.NewObjectID(),
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		Phone:          phone,
		Role:           role,
		HashedPassword: password,
		Verified:       notVerified,
	}

	err := validate.Struct(user)
	if err != nil {
		return User{}, err
	}

	err = userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == nil {
		return User{}, ErrAlreadyExists
	}

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUserByEmail(email string) (User, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	user := User{}

	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil

}
