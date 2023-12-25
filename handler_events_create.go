package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Bayan2019/mgrgServer/internal/auth"
	"github.com/Bayan2019/mgrgServer/internal/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID       string        `json:"id"`
	GuruID   string        `json:"pupil_id"`
	StartAt  time.Time     `json:"starts_on"`
	Duration time.Duration `json:"duration"`
	Essence  string        `json:"essence"`
	Location string        `json:"location"`
	Price    float32       `json:"price"`
}

func (cfg *apiConfig) handlerEventsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GuruID   string  `json:"event_id"`
		StartAt  string  `json:"start_at"`
		Duration string  `json:"duration"`
		Essence  string  `json:"essence"`
		Location string  `json:"location"`
		Price    float32 `json:"price"`
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

	guru_id, err := primitive.ObjectIDFromHex(subject)
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

	started_at, err := time.Parse(time.DateTime, params.StartAt)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse event's date")
		return
	}

	duration, err := time.ParseDuration(params.Duration)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse event's duration")
		return
	}

	event, err := db.CreateEvent(guru_id, started_at, duration,
		params.Essence, params.Location, params.Price)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create order")
		return
	}

	respondWithJSON(w, http.StatusCreated, Event{
		ID:       event.ID.Hex(),
		GuruID:   event.GuruID.Hex(),
		StartAt:  started_at,
		Duration: duration,
		Essence:  event.Essence,
		Location: event.Location,
		Price:    event.Price,
	})
}
