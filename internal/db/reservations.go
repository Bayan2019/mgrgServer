package db

import (
	"context"
	"fmt"

	// "net/http"
	"time"

	// "github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Reservation struct {
	ID      primitive.ObjectID `bson:"_id"`
	PupilID primitive.ObjectID `bson:"pupil_id"`
	EventID primitive.ObjectID `json:"event_id"`
	Price   float32            `json:"price"`
}

var validate = validator.New()
var reservationCollection *mongo.Collection = OpenCollection(Client, "reservations")

// var ErrNotExist = errors.New("resource does not exist")

func CreateReservation(pupil_id primitive.ObjectID,
	event_id primitive.ObjectID, price float32) (Reservation, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	pupil := User{}

	err := userCollection.FindOne(ctx, bson.M{"_id": pupil_id}).Decode(pupil)
	if err != nil {
		return Reservation{}, ErrNotExist
	}

	// if pupil.Role != "Pupil" {
	// 	return Reservation{}, errors.New("You are not a student to ")
	// }

	event := Event{}

	err = eventCollection.FindOne(ctx, bson.M{"_id": event_id}).Decode(event)
	if err != nil {
		return Reservation{}, ErrNotExist
	}

	reservation := Reservation{}

	err = reservationCollection.FindOne(ctx, bson.M{"pupil_id": pupil.ID, "event_id": event.ID}).Decode(&reservation)
	if err == nil {
		return Reservation{}, ErrAlreadyExists
	}

	reservation = Reservation{
		PupilID: pupil.ID,
		EventID: event.ID,
		Price:   price,
	}

	err = validate.Struct(reservation)
	if err != nil {
		return Reservation{}, err
	}

	_, err = reservationCollection.InsertOne(ctx, reservation)
	if err != nil {
		return Reservation{}, err
	}

	return reservation, nil
}

func GetReservationsByClassID(event_id primitive.ObjectID) ([]bson.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var reservations []bson.M

	cursor, err := reservationCollection.Find(ctx, bson.M{"event_id": event_id})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &reservations)
	if err != nil {
		return nil, err
	}
	defer cancel()
	fmt.Println(reservations)

	return reservations, nil
}
