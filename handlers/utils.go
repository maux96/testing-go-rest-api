package handlers

import (
	"errors"
	"my_rest_api/models"
	"my_rest_api/server"
	"net/http"
	"strings"

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
