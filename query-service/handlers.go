package main

import (
	"context"
	"github.com/digkill/giggle/db"
	"github.com/digkill/giggle/event"
	"github.com/digkill/giggle/schema"
	"github.com/digkill/giggle/search"
	"github.com/digkill/giggle/util"

	"log"
	"net/http"
	"strconv"
)

func onGiggleCreated(m event.GiggleCreatedMessage) {
	giggle := schema.Giggle{
		ID:        m.ID,
		Body:      m.Body,
		CreatedAt: m.CreatedAt,
	}
	if err := search.InsertGiggle(context.Background(), giggle); err != nil {
		log.Println(err)
	}
}

func listGigglesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	// Read parameters
	skip := uint64(0)
	skipStr := r.FormValue("skip")
	take := uint64(100)
	takeStr := r.FormValue("take")
	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
			return
		}
	}

	giggles, err := db.ListGiggles(ctx, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch giggles")
		return
	}

	util.ResponseOk(w, giggles)
}

func searchGigglesHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	// Read parameters
	query := r.FormValue("query")
	if len(query) == 0 {
		util.ResponseError(w, http.StatusBadRequest, "Missing query parameter")
		return
	}
	skip := uint64(0)
	skipStr := r.FormValue("skip")
	take := uint64(100)
	takeStr := r.FormValue("take")
	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
			return
		}
	}

	giggles, err := search.SearchGiggles(ctx, query, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseOk(w, []schema.Giggle{})
		return
	}

	util.ResponseOk(w, giggles)
}
