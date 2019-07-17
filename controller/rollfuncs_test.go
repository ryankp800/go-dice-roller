package controller

import (
	"log"
	"testing"
)

func TestRollDie(t *testing.T) {
	die := Dice{DiceValue: 20, Rolled: false}

	for i := 0; i < 1000; i++ {
		RollDie(&die)

		if die.Rolled != true {
			t.Errorf("Dice was not rolled, got: %t, want %t.", die.Rolled, true)
		}

		if die.RollValue < 0 || die.RollValue > 20 {
			t.Errorf("Dice roll value is outside acceptable range, got %d, expected between 1 and 20", die.RollValue)
		}
	}
}

func TestRoll(t *testing.T) {
	dieRoll := DiceRoll{
		DiceList: []Dice{
			{DiceValue: 20},
			{DiceValue: 12},
			{DiceValue: 10},
			{DiceValue: 8},
			{DiceValue: 6},
			{DiceValue: 4}},
		OverallRollValue: 0}

	for i := 0; i < 1000; i++ {
		Roll(&dieRoll)

		if dieRoll.OverallRollValue == 0 {
			t.Errorf("Expected value of >0 but was %v", dieRoll.OverallRollValue)
		}

		for _, die := range dieRoll.DiceList {
			if die.Rolled != true {
				t.Errorf("Dice was not rolled, got: %t, want %t.", die.Rolled, true)
			}

			if die.RollValue < 0 || die.RollValue > die.DiceValue {
				t.Errorf("Dice roll value is outside acceptable range, got %d, expected between 1 and %v", die.RollValue, die.DiceValue)
			}
		}
	}

}

func TestUpdateOrder(t *testing.T) {
	battle := Battle{
		Characters: []InitiativeRoll{
			{Name: "first", Modifier: 0, FinalValue: 20},
			{Name: "second", Modifier: 0, FinalValue: 19},
			{Name: "fifth", Modifier: 0, FinalValue: 4},
			{Name: "first", Modifier: 0, FinalValue: 18},
			{Name: "first", Modifier: 0, FinalValue: 17}},
	}

	battle = UpdateOrder(battle)

	// if battle.Characters[1].Order != 1 	|| battle.Characters[1].Order != 2 	|| battle.Characters[1].Order != 3 	|| battle.Characters[1].Order != 4 {
	// 	t.Errorf("Battle order was not updated %v", battle)
	// }

	for _, roll := range battle.Characters {
		log.Printf("finalValue %v, order %v", roll.FinalValue, roll.Order)
	}

}


func TestIncrimentOrder(t *testing.T) {
	battle := Battle{
		Characters: []InitiativeRoll{
			{Name: "first", Modifier: 0, FinalValue: 20, Order: 0},
			{Name: "second", Modifier: 0, FinalValue: 19, Order: 1},
			{Name: "third", Modifier: 0, FinalValue: 18, Order: 2},
			{Name: "fourth", Modifier: 0, FinalValue: 17, Order: 3}},
	}

	battle = IncrementOrder(battle)

	// if battle.Characters[1].Order != 1 	|| battle.Characters[1].Order != 2 	|| battle.Characters[1].Order != 3 	|| battle.Characters[1].Order != 4 {
	// 	t.Errorf("Battle order was not updated %v", battle)
	// }


}

