package server

import (
	"context"
	"errors"
	"log"
	"my_rest_api/database"
	"my_rest_api/repository"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

func (br *Broker) Config() *Config {
	return br.config
}

func NewServer(ctx context.Context, config *Config) (br *Broker, err error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("jwt secret is required")
	}
	if config.DatabaseUrl == "" {
		return nil, errors.New("database url is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}
	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)

	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	repository.SetRepository(repo)

	log.Println("Staring server on port", b.config.Port)
	if err := http.ListenAndServe("0.0.0.0:"+b.config.Port, b.router); err != nil {
		log.Fatalln("Server ListenAndServe", err.Error())
	}
}
