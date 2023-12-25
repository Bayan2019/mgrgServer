package db

import (
	"context"
	"errors"
	"time"

	// "github.com/golang/protobuf/ptypes/timestamp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Event struct {
	ID       primitive.ObjectID `json:"class_id"`
	GuruID   primitive.ObjectID `json:"guru_id"`
	StartAt  time.Time          `json:"starts_on"`
	Duration time.Duration      `json:"duration"`
	Essence  string             `json:"essence"`
	Location string             `json:"location"`
	Price    float32            `json:"price"`
}

var eventCollection *mongo.Collection = OpenCollection(Client, "events")

func CreateEvent(guru_id primitive.ObjectID, start_at time.Time, duration time.Duration,
	essence string, location string, price float32) (Event, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	guru := User{}

	err := userCollection.FindOne(ctx, bson.M{"_id": guru_id}).Decode(guru)
	if err != nil {
		return Event{}, ErrNotExist
	}

	if guru.Role != "Guru" {
		return Event{}, errors.New("you are not Guru to manage an event")
	}

	event := Event{
		GuruID:   guru.ID,
		StartAt:  start_at.UTC(),
		Duration: duration,
		Essence:  essence,
		Location: location,
		Price:    price,
	}

	err = validate.Struct(event)
	if err != nil {
		return Event{}, err
	}

	_, err = eventCollection.InsertOne(ctx, event)
	if err != nil {
		return Event{}, err
	}

	return event, nil
}
