package api

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
