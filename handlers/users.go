package handlers

import (
	"encoding/json"
	"my_rest_api/models"
	"my_rest_api/repository"
	"my_rest_api/server"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

const TOKEN_EXP_TIME_HOURS int = 1

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
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

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(requestObj.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := models.User{
			Id:             newId.String(),
			Email:          requestObj.Email,
			HashedPassword: string(hashedPass),
		}

		err = repository.InsertUser(r.Context(), &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignUpResponse{
			Id:    user.Id,
			Email: user.Email,
		})
	}
}

func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestObj SignUpRequest
		err := json.NewDecoder(r.Body).Decode(&requestObj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := repository.GetUserByEmail(r.Context(), requestObj.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user == nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(requestObj.Password)); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		claims := models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * time.Duration(TOKEN_EXP_TIME_HOURS)).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Token: tokenString,
		})

	}
}

func ProfileHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        claims, err := GetClaimsFromRequest(s, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

        user, err := repository.GetUserById(r.Context(), claims.UserId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")

        user.HashedPassword = ""
        json.NewEncoder(w).Encode(user)
	}
}
