package keypad

import (
	"machine"
)

type Driver struct {
	inputEnabled bool
	lastColumn   int //save position of the last column pressed
	lastRow      int //save position of the last row pressed
	columns      [4]machine.Pin
	rows         [4]machine.Pin
	mapping      [4][4]string
}

func (keypad *Driver) Configure(r4, r3, r2, r1, c4, c3, c2, c1 machine.Pin) {
	//The internal pullup resistor is going to hold the column to 5 V until a button is being pressed
	inputConfig := machine.PinConfig{Mode: machine.PinInputPullup}
	c4.Configure(inputConfig)
	c3.Configure(inputConfig)
	c2.Configure(inputConfig)
	c1.Configure(inputConfig)
	// Add the column pins to the columns array. By doing so, we can later just use a loop to iterate over all columns
	keypad.columns = [4]machine.Pin{c4, c3, c2, c1}

	outputConfig := machine.PinConfig{Mode: machine.PinOutput}
	r4.Configure(outputConfig)
	r3.Configure(outputConfig)
	r2.Configure(outputConfig)
	r1.Configure(outputConfig)
	// Add all the row pins to the rows array. This enables us to iterate over all the rows using a loop:
	keypad.rows = [4]machine.Pin{r4, r3, r2, r1}

	keypad.mapping = [4][4]string{
		{"1", "2", "3", "A"},
		{"4", "5", "6", "B"},
		{"7", "8", "9", "c"},
		{"*", "0", "#", "D"},
	}
	keypad.inputEnabled = true
	keypad.lastColumn = -1
	keypad.lastRow = -1
}

func (keypad *Driver) GetIndices() (int, int) {
	for rowIndex := range keypad.rows {

		rowPin := keypad.rows[rowIndex]
		rowPin.Low()

		for columnIndex := range keypad.columns {
			columnPin := keypad.columns[columnIndex]
			if !columnPin.Get() && keypad.inputEnabled {
				keypad.inputEnabled = false
				keypad.lastColumn = columnIndex
				keypad.lastRow = rowIndex
				rowPin.High()
				return keypad.lastRow, keypad.lastColumn
			}

			if columnPin.Get() &&
				columnIndex == keypad.lastColumn &&
				rowIndex == keypad.lastRow &&
				!keypad.inputEnabled {
				keypad.inputEnabled = true
			}

		}
		rowPin.High()

	}
	return -1, -1
}

func (keypad *Driver) GetGey() string {
	row, column := keypad.GetIndices()
	if row == -1 && column == -1 {
		return ""
	}

	return keypad.mapping[row][column]
}
