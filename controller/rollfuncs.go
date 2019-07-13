package controller

import (
	"math/rand"
)

//Roll Takes in a dice list and returns a new Rolled List
func Roll(roll *DiceRoll) {

	for i, die := range roll.DiceList {
		RollDie(&die)
		roll.OverallRollValue += die.RollValue
		roll.DiceList[i] = die
	}
}

//RollDie rolls and individual die
func RollDie(dice *Dice) {
	//Roll the dice
	dice.RollValue = rand.Intn(dice.DiceValue) + 1
	dice.Rolled = true
}
