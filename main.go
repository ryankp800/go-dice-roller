package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ryankp800/golang-dice-roller/api"
)

func main() {
	fmt.Println("Starting the application...")
	api.ConfigMongo()

	//TODO figure out what the host will be
	corsObj := handlers.AllowedOrigins([]string{"*"})

	r := mux.NewRouter()
	r.HandleFunc("/hello", api.HelloWorldHandler).Methods("GET")
	r.HandleFunc("/ws", api.HandleConnections)
	r.HandleFunc("/roll", api.RollDiceHandler)
	r.HandleFunc("/reset", api.ResetBattleHandler)
	r.HandleFunc("/getRoll/{id}", api.GetRollHandler).Methods("GET")
	r.Queries("value", "{value}")
	http.Handle("/", r)

	go api.HandleInitiative()

	// Apply the CORS middleware to our top-level router, with the defaults.
	http.ListenAndServe(GetPort(), handlers.CORS(corsObj)(r))
}

//GetPort Get the Port from the environment so we can run on Heroku
func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8000"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
