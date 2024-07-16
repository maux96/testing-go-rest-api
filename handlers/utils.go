package handlers

import (
	"errors"
	"my_rest_api/models"
	"my_rest_api/server"
	"net/http"
	"strings"
  "strconv"
	"github.com/golang-jwt/jwt"
)


func GetClaimsFromRequest(server server.Server, request *http.Request) (*models.AppClaims, error)  {
  tokenString := strings.TrimSpace(request.Header.Get("Authorization"))
  token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(t *jwt.Token) (interface{}, error) {
    return []byte(server.Config().JWTSecret), nil
  })
  if (err != nil) {
    return nil ,errors.New("problem in auth token")
  } 

  if appClaims,ok := token.Claims.(*models.AppClaims); ok && token.Valid {
    return appClaims, nil
  } else {
    return nil, errors.New("problem in auth token")
  }
}

func GetPaginationFromRequest(r *http.Request) (int64, int64, error) {
  var page, pageSize int64 = 0, 0

  pageAsString := r.URL.Query().Get("page")
  if pageAsString == "" {
    page = 0
  } else if p, err := strconv.ParseInt(pageAsString, 10, 64); err != nil {
    return 0, 0, err 
  } else  {
    page = p 
  }

  pageSizeAsString := r.URL.Query().Get("size")
  if pageSizeAsString == "" {
    pageSize = 8 
  } else if size, err := strconv.ParseInt(pageSizeAsString, 10, 64); err != nil {
    return 0, 0, err 
  } else  {
    pageSize = size 
  }
  
  return page, pageSize, nil
}
