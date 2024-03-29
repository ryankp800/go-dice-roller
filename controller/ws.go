package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var rollClients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Battle) // broadcast channel
var broadcastRolls = make(chan DiceResponse)

var currentBattle = Battle{Characters: []InitiativeRoll{}}

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

var HandleRollDiceConnection = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusForbidden)
	}
	log.Println(" connected to dice handler!")
	// ensure connection close when function returns
	defer ws.Close()
	rollClients[ws] = true

	ws.SetPingHandler(func(m string) error {
		go ws.WriteControl(websocket.PongMessage, []byte(m), time.Now().Add(time.Second*2))
		return nil
	})

	for {
		log.Println("Dice roller loop")
		var diceResponse DiceResponse
		var valueList RollRequest
		msgType, bytes, err := ws.ReadMessage()

		// We don't recognize any message that is not "ping".
		if msgType == websocket.TextMessage {
			log.Println("Received: ping.")
			continue
		} else {
			json.Unmarshal(bytes, &valueList)
		}

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
			broadCastDiceRoll(diceResponse)

		}
	}
})


func broadCastDiceRoll(response DiceResponse) {
	broadcastRolls <- response
}

// endTurn will end the current turn and increment order
func endTurn() {
	broadcast <- IncrementOrder(currentBattle)
}

func resetBattle() {
	currentBattle = Battle{InProgress: false, Characters: []InitiativeRoll{}}
	broadcast <- currentBattle
}

func startBattle() {
	currentBattle.InProgress = true
	broadcast <- currentBattle
}

func deleteFromBattle(id uuid.UUID) Battle {
	for i, char := range currentBattle.Characters {
		if char.ID == id {
			currentBattle.Characters = remove(currentBattle.Characters, i)
			broadcast <- currentBattle
			break
		}
	}
	return currentBattle
}

func overrideCharacterModifier(id uuid.UUID, mod int) {
	for i, char := range currentBattle.Characters {
		if char.ID == id {
			char.FinalValue, char.Modifier = (char.FinalValue-char.Modifier)+mod, mod
			currentBattle.Characters[i] = char
			currentBattle = UpdateOrder(currentBattle)
			broadcast <- currentBattle
			break
		}
	}
}

func remove(slice []InitiativeRoll, s int) []InitiativeRoll {
	return append(slice[:s], slice[s+1:]...)
}

// HandleInitConnection will handle incoming connections
var HandleInitConnection = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		log.Println("Running through loop")
		err := ws.ReadJSON(&initiativeRoll)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		} else {

			// If name is empty set username as name
			if initiativeRoll.Name == "" {
				initiativeRoll.Name = username
			}

			// For NPC set the owner as the user who inputs the roll
			initiativeRoll.Owner = username

			rollForInitiative(initiativeRoll, &currentBattle)

			broadcast <- currentBattle
		}
	}
})

// HandleInitiative websocket to handle the initiative rolls
func HandleInitiative() {
	for {
		msg := <-broadcast
		for client := range clients {
			log.Printf("clients %v", client)
			err := client.WriteJSON(msg)

			if err != nil {
				log.Printf("error: %v", err)
				err = client.Close()
				delete(clients, client)
			}
		}
	}
}

// BroadcastRoll sends new rolls to the channel and broadcasts to the websocket
func BroadcastRoll() {
	for {
		// grab next message from the broadcast channel
		msg := <-broadcastRolls
		// send it out to every client that is currently connected
		for client := range rollClients {
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
