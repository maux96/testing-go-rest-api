package handlers

import (
	"encoding/json"
	"my_rest_api/models"
	"my_rest_api/repository"
	"my_rest_api/server"
	"net/http"

	"github.com/segmentio/ksuid"
)

type SignUpRequest struct {
	Email string `json:"email"`
}

type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"string"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestObj SignUpRequest
		err := json.NewDecoder(r.Body).Decode(&requestObj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newId, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := models.User{
			Id:    newId.String(),
			Email: requestObj.Email,
		}

		err = repository.InsertUser(r.Context(), &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignUpResponse{
			Id:    user.Id,
			Email: user.Email,
		})
	}
}
