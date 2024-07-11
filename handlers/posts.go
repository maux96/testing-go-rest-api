package handlers

import (
	"encoding/json"
	"my_rest_api/models"
	"my_rest_api/repository"
	"my_rest_api/server"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type UpsertPostRequest struct {
	PostContent string `json:"post_content"`
}
type InsertPostResponse struct {
  Id          string `json:"id"`
  // UserId          string `json:"user_id"`
	PostContent string `json:"post_content"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
   	var requestObj UpsertPostRequest 
		err := json.NewDecoder(r.Body).Decode(&requestObj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

    if appClaims, err := GetClaimsFromRequest(s, r) ; err == nil {

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
      http.Error(w, err.Error(), http.StatusUnauthorized)
    }
   
  }
}

func GetPostByIdHandler(s server.Server) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    post, err := repository.GetPostById(r.Context(), params["id"])
    if (err != nil) {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
     
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
  } 
}

func UpdatePostHandler(s server.Server) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    var requestObj UpsertPostRequest
    err := json.NewDecoder(r.Body).Decode(&requestObj)
    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    } 

    appClaims, err := GetClaimsFromRequest(s, r)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    } 

    err = repository.UpdatePost(r.Context(), params["id"], &models.Post{PostContent: requestObj.PostContent, UserId: appClaims.UserId})
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    
  }
}


func DeletePostHandler(s server.Server) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    appClaims, err := GetClaimsFromRequest(s, r)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }


    err = repository.DeletePostById(r.Context(), params["id"], appClaims.UserId)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }
}

func ListPostHandlers(s server.Server) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    pageAsString := r.URL.Query().Get("page")
    var page int64 = 0 

    if p, err := strconv.ParseInt(pageAsString, 10, 64); pageAsString != "" && err == nil {
      page = int64(p)
    }

    posts, err := repository.ListPosts(r.Context(), page)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
  } 
}


