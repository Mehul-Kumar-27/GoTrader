package main

import (
	"gotrader/logger"
	"gotrader/server/cmd/api/socket"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	logger := logger.CreateCustomLogger("server/api")

	r := mux.NewRouter()
	m := socket.NewManger()
	r.HandleFunc("/api/stocks", m.InitalRoute).Methods("GET")
	r.HandleFunc("/api/socket", socket.HandleWebSockerConnection)

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	err := server.ListenAndServe()
	if err != nil {
		logger.Panicf("Failed to listen: %v", err)
	}
}
