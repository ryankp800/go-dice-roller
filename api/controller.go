package api

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//HelloWorldHandler hello world
func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(`{"hello": "there"}`)
}

//GetRollHandler gets a roll from the database based on the object ID
func GetRollHandler(w http.ResponseWriter,
	r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	roll := getDiceRollByID(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roll)

}
//RollDiceHandler roll dice yo
func RollDiceHandler(w http.ResponseWriter,
	r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Get the dice value list from the rul
	queryVariableMap := r.URL.Query()

	//Get the value property and put it into an array
	valueList := queryVariableMap["value"]

	//Setup regex to split on the 'd' value
	re := regexp.MustCompile("[Dd]")
	var dieList DiceRoll

	//First iterate through all of the value objects in the list
	for _, value := range valueList {

		//Use regex to split the d value
		me := re.Split(value, -1)

		//Check if there are greater than 0 dice
		if numOfDie, err := strconv.Atoi(me[0]); err == nil {

			//For each die of that vlaue being rolled create a die object
			for j := 0; j < numOfDie; j++ {

				//Convert the Value string to an int
				if dVal, err := strconv.Atoi(me[1]); err == nil {

					//Create die of the value and then appen it to the die list
					die := Dice{DiceValue: dVal, ID: 1, Rolled: false, RollValue: 0}
					dieList.DiceList = append(dieList.DiceList, die)
				} else {
					log.Println(err)
				}
			}
		}

	}

	//Roll the completed dielist
	dieList = Roll(dieList)
	dieList.ID = primitive.NewObjectID()

	insertDiceRoll(dieList)
	//Marshal the list
	data, err := json.Marshal(dieList)
	if err != nil {
		log.Fatalf("JSON marshalling failed: %s", err)
	}

	//Write to the response
	w.Write(data)

}

//Roll Takes in a dice list and returns a new Rolled List
func Roll(roll DiceRoll) DiceRoll {

	for i, dice := range roll.DiceList {
		//Roll the dice
		rollValue := rand.Intn(dice.DiceValue) + 1
		roll.DiceList[i].Rolled = true
		roll.OverallRollValue += rollValue

		//Set the value
		roll.DiceList[i].RollValue = rollValue
	}
	return roll
}
