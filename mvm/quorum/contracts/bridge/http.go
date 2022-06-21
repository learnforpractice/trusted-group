package main

// FIXME do rate limit based on IP
// POST /users
// GET /users/:id

import (
	"encoding/json"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/unrolled/render"
)

var (
	proxy *Proxy
	store *Storage
)

func StartHTTP(p *Proxy, s *Storage) error {
	proxy, store = p, s
	router := httptreemux.New()
	router.POST("/users", createUser)
	router.GET("/users/:id", readUser)
	return http.ListenAndServe(":3000", router)
}

func readUser(w http.ResponseWriter, r *http.Request, params map[string]string) {
	user, err := store.readUser(params["id"])
	if err != nil {
		render.New().JSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err})
	} else if user == nil {
		render.New().JSON(w, http.StatusNotFound, map[string]interface{}{})
	} else {
		render.New().JSON(w, http.StatusOK, map[string]interface{}{"user": user})
	}
}

func createUser(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body struct {
		PublicKey string `json:"public_key"`
		Signature string `json:"signature"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		render.New().JSON(w, http.StatusBadRequest, map[string]interface{}{"error": err})
		return
	}
	user, err := proxy.createUser(r.Context(), store, body.PublicKey, body.Signature)
	if err != nil {
		render.New().JSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err})
		return
	}
	render.New().JSON(w, http.StatusOK, map[string]interface{}{"user": user})
}
