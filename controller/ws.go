package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Battle)            // broadcast channel
var BroadcastRolls = make(chan DiceResponse)

var currentBattle Battle

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleConnections will handle incoming connections
var HandleConnections = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusForbidden)
	}
	log.Println("connected!")
	// ensure connection close when function returns
	defer ws.Close()
	clients[ws] = true

	// user := r.Context().Value("user")
	// k, _ := user.(*jwt.Token).Claims.(jwt.MapClaims)
	// username := k["username"].(string)
	username := "joe"
	for {
		var initiativeRoll InitiativeRoll
		err := ws.ReadJSON(&initiativeRoll)
		// Read in a new message as JSON and map it to a Message object

		// If name is empty set username as name
		if initiativeRoll.Name == "" {
			initiativeRoll.Name = username
		}
		// For NPC set the owner as the user who inputs the roll
		initiativeRoll.Owner = username

		rollForInitiative(initiativeRoll, &currentBattle)

		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)

		}
		// send the new message to the broadcast channel
		broadcast <- currentBattle
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

// HandleInitiative websocket to handle the initiative rolls
func HandleInitiative() {
	for {
		// grab next message from the broadcast channel
		msg := <-broadcast
		// send it out to every client that is currently connected
		for client := range clients {
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

// BroadcastRoll sends new rolls to the channel and broadcasts to the websocket
func BroadcastRoll() {
	for {
		// grab next message from the broadcast channel
		msg := <-BroadcastRolls
		// send it out to every client that is currently connected
		for client := range clients {
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
