package services

import (
	"fmt"
	"net/http"
	"websocket-demo/internal/models"
	"websocket-demo/internal/utils"
)

type Server struct {
	config *models.Config
}

func NewServer(config *models.Config) *Server {
	return &Server{
		config: config,
	}
}

func (this *Server) Start() {
	http.HandleFunc("/health", health)
	err := http.ListenAndServe(fmt.Sprintf(":%d", this.config.RestPort), nil)
	utils.PanicOnError(err)
}

func health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Health: Ok!")
}
