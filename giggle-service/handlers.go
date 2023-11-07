package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/digkill/giggle/db"
	"github.com/digkill/giggle/event"
	"github.com/digkill/giggle/schema"
	"github.com/digkill/giggle/util"
	"github.com/segmentio/ksuid"
)

func createGiggleHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID string `json:"id"`
	}

	ctx := r.Context()

	// Read parameters
	body := template.HTMLEscapeString(r.FormValue("body"))
	if len(body) < 1 || len(body) > 140 {
		util.ResponseError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	// Create giggle
	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandomWithTime(createdAt)
	if err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create giggle")
		return
	}
	giggle := schema.Giggle{
		ID:        id.String(),
		Body:      body,
		CreatedAt: createdAt,
	}
	if err := db.InsertGiggle(ctx, giggle); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create giggle")
		return
	}

	// Publish event
	if err := event.PublishGiggleCreated(giggle); err != nil {
		log.Println(err)
	}

	util.ResponseOk(w, response{ID: giggle.ID})
}
