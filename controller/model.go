package controller

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Dice a basic Dice object that carries an ID a value a roll value and a rolled flag
type Dice struct {
	ID        int  `json:"id,omitempty" bson:"_id,omitempty"`
	DiceValue int  `json:"dValue,omitempty"`
	Rolled    bool `json:"rolled,omitempty"`
	RollValue int  `json:"rollValue,omitempty"`
}

//DiceRoll a list of Dice that contains an overallRollValue
type DiceRoll struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	DiceList         []Dice             `json:"diceList,omitempty"`
	OverallRollValue int                `json:"overallRollResult,omitempty"`
}

//InitiativeRoll a single PC initiative roll within a battle
type InitiativeRoll struct {
	Name            string `json:"name,omitempty"`
	PlayerCharacter bool   `json:"player_character,omitempty"`
	Advantage       bool   `json:"advantage,omitempty"`
	Modifier        int    `json:"modifier,omitempty"`
	FinalValue      int    `json:"final_value,omitempty"`
}

//Battle object that structures the current fight
type Battle struct {
	Characters []InitiativeRoll `json:"participants,omitempty"`
	IsComplete bool             `json:"is_complete,omitempty"`
}

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Token     string `json:"token,omitempty"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

type DiceResponse struct {
	DiceRoll DiceRoll  `json:"dice_roll"`
	User User `json:"user"`
}