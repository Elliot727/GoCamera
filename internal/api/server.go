package api

import (
	"fmt"
	"net/http"
)

type Server struct {
	Port string
}

func NewServer(port string) *Server {
	return &Server{Port: port}
}

func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Welcome to the API")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/", s.RootHandler)

	fmt.Printf("Starting server on port %s...\n", s.Port)
	return http.ListenAndServe(":"+s.Port, nil)
}
