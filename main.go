package main

import (
	"encoding/json"
	"io"
	"log"
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
	http.ListenAndServe(":8000", handlers.CORS(corsObj)(r))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, `{"Hello": "World"}`)
}

func rollDiceHandler(w http.ResponseWriter,
	r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	number := r.URL.Query()
	log.Println(number)
	valueList := number["value"]

	log.Println(valueList)
	re := regexp.MustCompile("[Dd]")
	var dieList DiceRoll
	for i, value := range valueList {
		me := re.Split(value, -1)
		log.Println("value {i} {s}", i, me)
		log.Println("value [0] then [1]", me[0], me[1])

		if numOfDie, err := strconv.Atoi(me[0]); err == nil {

			for j := 0; j < numOfDie; j++ {
				if dVal, err := strconv.Atoi(me[1]); err == nil {
					die := Dice{DiceValue: dVal, ID: 1, Rolled: false, RollValue: 0}
					dieList.DiceList = append(dieList.DiceList, die)
				} else {
					log.Println(err)
				}
			}
		}

	}

	dieList = Roll(dieList)

	log.Println(dieList)

	data, err := json.Marshal(dieList)
	if err != nil {
		log.Fatalf("JSON marshalling failed: %s", err)
	}

	w.Write(data)

}

//Roll Takes in a dice list and returns a new Rolled ListS
func Roll(roll DiceRoll) DiceRoll {

	for i, dice := range roll.DiceList {
		rollValue := rand.Intn(dice.DiceValue) + 1
		log.Println("Dice", dice.RollValue)
		roll.DiceList[i].Rolled = true
		roll.OverallRollValue += rollValue
		log.Println("Dice", dice)
		roll.DiceList[i].RollValue = rollValue
	}

	log.Println(roll)

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
