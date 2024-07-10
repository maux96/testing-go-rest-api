package handlers

import (
	"encoding/json"
	"my_rest_api/models"
	"my_rest_api/repository"
	"my_rest_api/server"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
)

type InsertPostRequest struct {
	PostContent string `json:"post_content"`
}
type InsertPostResponse struct {
  Id          string `json:"id"`
  // UserId          string `json:"user_id"`
	PostContent string `json:"post_content"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
   	var requestObj InsertPostRequest 
		err := json.NewDecoder(r.Body).Decode(&requestObj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

    tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
    token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(t *jwt.Token) (interface{}, error) {
      return []byte(s.Config().JWTSecret), nil
    })

    if appClaims,ok := token.Claims.(*models.AppClaims); ok && token.Valid {

		  newId, err := ksuid.NewRandom()
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
      }

      post := models.Post{
        Id:             newId.String(),
        PostContent:    requestObj.PostContent,
        UserId:         appClaims.UserId,
      }

      err = repository.InsertPost(r.Context(), &post)
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
      }

      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(InsertPostResponse {
        Id:           post.Id,
        PostContent:  post.PostContent,
      })
    } else {
      http.Error(w, "", http.StatusUnauthorized)
    }
   
  }
}
