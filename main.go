package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	corsObj := handlers.AllowedOrigins([]string{"*"})

	r := mux.NewRouter()
	r.HandleFunc("/hi", helloWorldHandler)
	r.HandleFunc("/roll", rollDiceHandler)
	r.Queries("value", "{value}")
	http.Handle("/", r)
	// Apply the CORS middleware to our top-level router, with the defaults.
	http.ListenAndServe(GetPort(), handlers.CORS(corsObj)(r))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, `{"Hello": "World"}`)
}

func rollDiceHandler(w http.ResponseWriter,
	r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	//Get the dice value list from the rul
	number := r.URL.Query()

	//Get the value property and put it into an array
	valueList := number["value"]

	//Setup regex to split on the 'd' value
	re := regexp.MustCompile("[Dd]")
	var dieList DiceRoll

	//First iterate through all of the value objects in the list
	for i, value := range valueList {

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

	//Marshal the list
	data, err := json.Marshal(dieList)
	if err != nil {
		log.Fatalf("JSON marshalling failed: %s", err)
	}

	//Write to the response
	w.Write(data)

}

//Roll Takes in a dice list and returns a new Rolled ListS
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

//Dice a basic Dice object that carries an ID a value a roll value and a rolled flag
type Dice struct {
	ID        int  `json:"id"`
	DiceValue int  `json:"dValue"`
	Rolled    bool `json:"rolled"`
	RollValue int  `json:"rollValue"`
}

//DiceRoll a list of Dice that contains an overallRollValue
type DiceRoll struct {
	DiceList         []Dice `json:"diceList"`
	OverallRollValue int    `json:"overallRollResult"`
}

// Get the Port from the environment so we can run on Heroku
func GetPort() string {
	 	var port = os.Getenv("PORT")
	 	// Set a default port if there is nothing in the environment
	 	if port == "" {
	 		port = "8000"
	 		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	 	}
	 	return ":" + port
	 }
