package controller

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var rollClients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Battle)            // broadcast channel
var BroadcastRolls = make(chan DiceResponse)

var currentBattle  = Battle{Characters: []InitiativeRoll{}}

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}



func extractClaims(tokenString string) string {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		panic(err)
	}
	// ... error handling
	// do something with decoded claims
	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val)
		if key == "username" && val != "" {
			return fmt.Sprintf("%v", val)
		}

	}

	return "no username"
}

var HandleConnections = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusForbidden)
	}
	log.Println(" connected to dice handler!")
	// ensure connection close when function returns
	defer ws.Close()
	rollClients[ws] = true

	for {
		log.Println("Dice roller loop")
		var diceResponse DiceResponse
		var valueList RollRequest
		err := ws.ReadJSON(&valueList)
		if err != nil {
			log.Println("Dice roller loop ")
			delete(rollClients, ws)
			break
		} else {
			re := regexp.MustCompile("[Dd]")
			dieList := extractDieList(valueList.ValString, re)

			// Roll the completed dielist
			Roll(&dieList)
			dieList.ID = primitive.NewObjectID()
			insertDiceRoll(dieList)

			diceResponse.User.Username, diceResponse.DiceRoll = "fake", dieList

			// send the new message to the broadcast channel
			BroadcastRolls <- diceResponse

		}
	}
})

// StartBattle handles the initial set up of a new encounter
func StartBattle(r *http.Request, ws *websocket.Conn) {


}

// endTurn will end the current turn and increment order
func endTurn() {
	broadcast <- IncrementOrder(currentBattle)
}

func resetBattle() {
	currentBattle = Battle{InProgress:false}
}

func startBattle() {
	currentBattle.InProgress = true
}
// HandleInitConnections will handle incoming connections
var HandleInitConnections = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	// Get the token from the url
	queryVariableMap := r.URL.Query()
	tokenString := queryVariableMap["token"][0]
	username := extractClaims(tokenString)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusForbidden)
	}
	log.Println("connected!")
	// ensure connection close when function returns
	defer ws.Close()

	clients[ws] = true

	// Push existing data out on connection
	broadcast <- currentBattle

	var initiativeRoll InitiativeRoll
	for {
		fmt.Println("Running through loop")
		err := ws.ReadJSON(&initiativeRoll)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		} else {
			// Read in a new message as JSON and map it to a Message object

			// If name is empty set username as name
			if initiativeRoll.Name == "" {
				initiativeRoll.Name = username
			}

			// For NPC set the owner as the user who inputs the roll
			// initiativeRoll.Owner = username
			rollForInitiative(initiativeRoll, &currentBattle)

			broadcast <- currentBattle
		}
	}
})

// HandleInitiative websocket to handle the initiative rolls
func HandleInitiative() {
	log.Println("HandleInitiative() started")
	for {
		log.Println("HandleInitiative() looped")
		// grab next message from the broadcast channel
		msg := <-broadcast
		// send it out to every client that is currently connected
		log.Printf("clijents: %v", clients)
		for  client := range clients {
			// log.Println("handling init --> sending to client %v", client)
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v, client # %v", err)
				// TODO I dont think we want to close the client in this situation
				err = client.Close()
				delete(clients, client)
			}
		}
	}
}

// BroadcastRoll sends new rolls to the channel and broadcasts to the websocket
func BroadcastRoll() {
	log.Println("BroadcastRoll() started")
	for {
		log.Println("BroadcastRoll() looped")
		// grab next message from the broadcast channel
		msg := <-BroadcastRolls
		// send it out to every client that is currently connected
		for client := range rollClients {
			log.Println("meh")
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				// TODO I dont think we want to close the client in this situation
				client.Close()
				delete(clients, client)
			}
		}
	}

}
