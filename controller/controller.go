package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"

	guuid "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HelloWorldHandler hello world
func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(`{"hello": "there"}`)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
	}
}



// GetRollHandler gets a roll from the database based on the object ID
func GetRollHandler(w http.ResponseWriter,
	r *http.Request) {
	setupResponse(&w, r)
	params := mux.Vars(r)
	id := params["id"]

	roll := GetDiceRollByID(id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roll)

}

var RollDiceHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	setupResponse(&w, r)
	var diceResponse DiceResponse
	user := r.Context().Value("user")
	k, _ := user.(*jwt.Token).Claims.(jwt.MapClaims)
			diceResponse.User.Username = k["username"].(string)


	// Get the dice value list from the rul
	queryVariableMap := r.URL.Query()

	// Get the value property and put it into an array
	valueList := queryVariableMap["value"]

	// Setup regex to split on the 'd' value
	re := regexp.MustCompile("[Dd]")
	dieList := extractDieList(valueList, re)

	// Roll the completed dielist
	Roll(&dieList)
	dieList.ID = primitive.NewObjectID()

	insertDiceRoll(dieList)

	diceResponse.DiceRoll = dieList

	// Marshal the list
	data, err := json.Marshal(diceResponse)
	if err != nil {
		log.Printf("JSON marshalling failed: %s", err)
		return
	}

	// Send dice roll to channel to be sent over ws
	BroadcastRolls <- diceResponse


	// Write to the response
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not parse die data %s", err)
		return
	}
})

func getusername(tokenString string, diceResponse *DiceResponse) {

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		diceResponse.User.Username = claims["username"].(string)

	} else {
		diceResponse.User.Username = "undefined"
	}
}


func extractDieList(valueList []string, re *regexp.Regexp) DiceRoll {
	var dieList DiceRoll
	// First iterate through all of the value objects in the list
	for _, value := range valueList {

		// Use regex to split the d value
		me := re.Split(value, -1)

		// Check if there are greater than 0 dice
		if numOfDie, err := strconv.Atoi(me[0]); err == nil {

			// For each die of that value being rolled create a die object
			for j := 0; j < numOfDie; j++ {

				// Convert the Value string to an int
				if dVal, err := strconv.Atoi(me[1]); err == nil {

					// Create die of the value and then appen it to the die list
					die := Dice{DiceValue: dVal, ID: 1, Rolled: false, RollValue: 0}
					dieList.DiceList = append(dieList.DiceList, die)
				} else {
					log.Println(err)
				}
			}
		}

	}
	return dieList
}

var DeleteInitHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	params := mux.Vars(r)
	id := params["id"]

	uuid, err := guuid.Parse(id); if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	battle := deleteFromBattle(uuid)

	b, _ := json.Marshal(battle)
	w.Write(b)

	return
})

var UpdateModifierHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	params := mux.Vars(r)
	id := params["id"]

	b, err := ioutil.ReadAll(r.Body); 	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var mod Modifier
	err = json.Unmarshal(b, &mod); if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error unmarshalling initiative list err %v", err)
		return
	}

	uuid, err := guuid.Parse(id)
	overrideCharacterModifier(uuid, mod.Mod)

	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(currentBattle)
	w.Write(data)

	return
})

var InitiativeHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	user := r.Context().Value("user")
	k, _ := user.(*jwt.Token).Claims.(jwt.MapClaims)
	username := k["username"].(string)


	var initRollList InitiativeRollList

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(b, &initRollList); if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error unmarshalling initiative list err %v", err)
		return
	}


	for _, initRoll := range initRollList.CharacterList {
		// If name is empty set username as name
		if initRoll.Name == "" {
			initRoll.Name = username
		}

		// For NPC set the owner as the user who inputs the roll
		initRoll.Owner = username


		rollForInitiative(initRoll, &currentBattle)



		broadcast <- currentBattle

	}

	w.WriteHeader(http.StatusCreated)
	battleBytes, _ := json.Marshal(currentBattle)
	w.Write(battleBytes)

	return
})

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	var user User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	var res ResponseResult
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	var result User
	err = UserCollection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				res.Error = "Error While Hashing Password, Try Again"
				json.NewEncoder(w).Encode(res)
				return
			}
			user.Password = string(hash)

			_, err =  UserCollection.InsertOne(context.TODO(), user)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				res.Error = "Error While Creating User, Try Again"
				json.NewEncoder(w).Encode(res)
				return
			}
			res.Result = "Registration Successful"
			json.NewEncoder(w).Encode(res)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Result = "Username already Exists!!"
	json.NewEncoder(w).Encode(res)
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	var user User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusForbidden)
	}

	var result User
	var res ResponseResult

	err = UserCollection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		res.Error = "Invalid username"
		json.NewEncoder(w).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		res.Error = "Invalid password"
		json.NewEncoder(w).Encode(res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": result.Username,
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		res.Error = "Error while generating token,Try again"
		json.NewEncoder(w).Encode(res)
		return
	}

	result.Token = tokenString
	result.Password = ""

	json.NewEncoder(w).Encode(result)

}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	var result User
	var res ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Username = claims["username"].(string)

		json.NewEncoder(w).Encode(result)
		return
	} else {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

}

var EndTurnHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	setupResponse(&w, r)
	endTurn()
	json.NewEncoder(w).Encode(`{"turnEnded": true}`)
})

var StartBattleHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	setupResponse(&w, r)
	startBattle()
	json.NewEncoder(w).Encode(`{"battleStarted": true}`)
})

// ResetBattleHandler clears battle object
var ResetBattleHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	resetBattle()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(`{"battleReset": true}`)
})

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
