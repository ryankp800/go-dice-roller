package controller

import (
	"fmt"
	"math/rand"
	"sort"
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

// ByFinalValue implements sort.Interface for []Person based on
// the Age field.
type ByFinalValue []InitiativeRoll

func (a ByFinalValue) Len() int           { return len(a) }
func (a ByFinalValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFinalValue) Less(i, j int) bool { return a[i].FinalValue < a[j].FinalValue }

// ByFinalValue implements sort.Interface for []Person based on
// the Age field.
type ByOrderValue []InitiativeRoll

func (a ByOrderValue) Len() int           { return len(a) }
func (a ByOrderValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOrderValue) Less(i, j int) bool { return a[i].Order < a[j].Order }

//
func UpdateOrder(battle Battle) Battle {

	sort.Sort(ByFinalValue(battle.Characters))

	for i, _ := range battle.Characters {
		battle.Characters[i].Order = i
	}
	fmt.Println(battle)

	return battle
}

func IncrementOrder(battle Battle) Battle {

	sort.Sort(ByOrderValue(battle.Characters))

	for i, k := range battle.Characters {
		if k.Order != 0 {
			battle.Characters[i].Order = i-1
		} else {
			battle.Characters[i].Order = len(battle.Characters) -1
		}
	}

	sort.Sort(ByOrderValue(battle.Characters))
	fmt.Println(battle)

	return battle
}

