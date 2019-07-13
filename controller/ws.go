package controller

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
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
		log.Fatal(err)
	}
	log.Println("connected!")
	// ensure connection close when function returns
	defer ws.Close()
	clients[ws] = true

	user := r.Context().Value("user")
	k, _ := user.(*jwt.Token).Claims.(jwt.MapClaims)
	username := k["username"].(string)


	for {
		var initiativeRoll InitiativeRoll
		initiativeRoll.Name = username
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&initiativeRoll)
		rollForInitiative(initiativeRoll, &currentBattle)

		log.Println(initiativeRoll)

		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// send the new message to the broadcast channel
		log.Println(currentBattle)
		broadcast <- currentBattle
	}
})

func rollForInitiative(roll InitiativeRoll, battle *Battle) {

	die := Dice{DiceValue: 20}
	RollDie(&die)

	roll.FinalValue = die.RollValue + roll.Modifier

	battle.Characters = append(battle.Characters, roll)
}

// HandleInitiative websocket to handle the initiative rolls
func HandleInitiative() {
	for {
		// grab next message from the broadcast channel
		msg := <-broadcast
		// send it out to every client that is currently connected
		for client := range clients {
			log.Println("my battle", msg)
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func resetBattle() {
	currentBattle = Battle{}
}

func BroadcastRoll() {
	for {
		// grab next message from the broadcast channel
		msg := <-BroadcastRolls
		// send it out to every client that is currently connected
		for client := range clients {
			log.Println("my roll", msg)
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}

}
