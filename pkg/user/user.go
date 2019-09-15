package user

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type User struct {
	ID int `json:"user_id"`
}

// POST /api/v1/user/update {"id": 1, .... }
func Update(r *http.Request) (*http.Response, error) {
	// TODO
	return &http.Response{}, nil
}


// GET /api/v1/user/delete?ids=[1]
func Delete(r *http.Request) (*http.Response, error) {
	// TODO
	return &http.Response{}, nil
}


func ResponseJSON (w http.ResponseWriter, v interface{}) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return
	}
}


func Get(w http.ResponseWriter, r *http.Request) {
	log.Info("Get request to user")
	ResponseJSON(w, &User{ID: 1})
}

func Add(w http.ResponseWriter, r *http.Request) {
	log.Info("Get request to user/add")
	ResponseJSON(w, &User{ID: 2})
}
