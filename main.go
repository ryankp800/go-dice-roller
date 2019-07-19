package main

import (
	"fmt"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/ryankp800/golang-dice-roller/controller"
)

func main() {
	fmt.Println("Starting the application...")
	controller.ConfigMongo()


	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	r := routes()

	go controller.HandleInitiative()
	go controller.BroadcastRoll()

	// Apply the CORS middleware to our top-level router, with the defaults.
	http.ListenAndServe(getPort(), c.Handler(r))
}

func routes() *mux.Router {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})
	r := mux.NewRouter()


	r.HandleFunc("/health", controller.HelloWorldHandler).Methods("GET")
	r.Handle("/ws/init", controller.HandleInitConnections)
	r.Handle("/ws/roll", controller.HandleConnections)
	r.Handle("/roll", jwtMiddleware.Handler(controller.RollDiceHandler))
	r.Handle("/end", jwtMiddleware.Handler(controller.EndTurnHandler)).Methods("GET")
	r.HandleFunc("/reset", controller.ResetBattleHandler)
	r.HandleFunc("/getRoll/{id}", controller.GetRollHandler).Methods("GET")
	r.HandleFunc("/register", controller.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", controller.LoginHandler).Methods("POST")
	r.HandleFunc("/profile", controller.ProfileHandler).Methods("GET")
	r.Queries("value", "{value}")
	http.Handle("/", r)
	return r
}

func getPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8000"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}