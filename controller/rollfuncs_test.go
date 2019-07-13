package controller

import (
	"testing"
)

func TestRollDie(t *testing.T) {
	die := Dice{DiceValue: 20, Rolled: false}

	for i := 0; i < 1000; i ++ {
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
			{DiceValue:20},
			{DiceValue:12},
			{DiceValue:10},
			{DiceValue:8},
			{DiceValue:6},
			{DiceValue:4}},
	OverallRollValue:0}

	for i :=0; i < 1000 ; i++  {
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


