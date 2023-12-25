package main

import (
	"encoding/json"
	"net/http"

	"github.com/Bayan2019/mgrgServer/internal/auth"
	"github.com/Bayan2019/mgrgServer/internal/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reservation struct {
	ID      string `json:"id"`
	PupilID string `json:"pupil_id"`
	EventID string `json:"event_id"`
}

func (cfg *apiConfig) handlerReservationsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		EventID string  `json:"event_id"`
		Price   float32 `json:"price"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	pupil_id, err := primitive.ObjectIDFromHex(subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse pupil ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	event_id, err := primitive.ObjectIDFromHex(params.EventID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse class ID")
		return
	}

	order, err := db.CreateReservation(pupil_id, event_id, params.Price)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create order")
		return
	}

	respondWithJSON(w, http.StatusCreated, Reservation{
		ID:      order.ID.Hex(),
		PupilID: order.PupilID.Hex(),
		EventID: order.EventID.Hex(),
	})
}
