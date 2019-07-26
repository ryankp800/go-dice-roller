package controller

import (
	"math/rand"
	"sort"
	"time"

	guuid "github.com/google/uuid"
)

var seeded = false

// Roll Takes in a dice list and returns a new Rolled List
func Roll(roll *DiceRoll) {

	for i, die := range roll.DiceList {
		RollDie(&die)
		roll.OverallRollValue += die.RollValue
		roll.DiceList[i] = die
	}
}

// RollDie rolls and individual die
func RollDie(dice *Dice) {
	if !seeded {
		rand.Seed(time.Now().UTC().UnixNano())
		seeded = true
	}
	dice.RollValue = rand.Intn(dice.DiceValue) + 1
	dice.Rolled = true
}

// UpdateOrder takes in a battle and will sort the characters based off of their final roll value
func UpdateOrder(battle Battle) Battle {

	// TODO this only works at the beginning of a battle. Need to decide how to handle enemies entering a battle

	sort.Sort(ByFinalValue(battle.Characters))

	for i := range battle.Characters {
		battle.Characters[i].Order = i
	}

	return battle
}

func rollForInitiative(roll InitiativeRoll, battle *Battle) {

	die := Dice{DiceValue: 20}
	RollDie(&die)
	if roll.Advantage {
		advDie := Dice{DiceValue: 20}
		RollDie(&advDie)

		if advDie.RollValue > die.RollValue {
			die = advDie
		}
	}

	roll.ID, roll.FinalValue = guuid.New(), die.RollValue+roll.Modifier

	battle.Characters = append(battle.Characters, roll)

	UpdateOrder(*battle)
}

func IncrementOrder(battle Battle) Battle {

	sort.Sort(ByOrderValue(battle.Characters))

	for i, k := range battle.Characters {
		if k.Order != 0 {
			battle.Characters[i].Order = i - 1
		} else {
			battle.Characters[i].Order = len(battle.Characters) - 1
		}
	}

	sort.Sort(ByOrderValue(battle.Characters))
	// fmt.Println(battle)

	return battle
}

// ByFinalValue implements sort.Interface for []InitiativeRoll based on
// the FinalValue field, sorting highest to lowest
type ByFinalValue []InitiativeRoll

func (a ByFinalValue) Len() int           { return len(a) }
func (a ByFinalValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFinalValue) Less(i, j int) bool { return a[i].FinalValue > a[j].FinalValue }

// ByOrderValue implements sort.Interface for []InitiativeRoll based on
// the Order field.
type ByOrderValue []InitiativeRoll

func (a ByOrderValue) Len() int           { return len(a) }
func (a ByOrderValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOrderValue) Less(i, j int) bool { return a[i].Order < a[j].Order }
