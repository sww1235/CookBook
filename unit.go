package main

import "fmt"

type Unit struct {
	ID          int    // id of unit in database
	Name        string // human readable name of unit
	Symbol      string // recipe symbol
	Description string // unit description
	IsCustom    bool   // is unit custom or standard

}

func (u Unit) String() string {
	stringString := ""
	stringString += fmt.Sprintf("Unit Name: %s\nUnit Symbol: %s\n", u.Name, u.Symbol)
	return stringString

}

// Type conversion is a storage method for conversions between different units
// The conversions use the following formula.
// toUnitValue = ((fromUnitValue + fromOffset) * multiplicand / denominator) + toOffset
type conversion struct {
	ID           int     // id of conversion factor in database
	FromUnit     int     // db id of unit to convert from
	ToUnit       int     // db id of unit to convert to
	multiplicand float64 //
	denominator  float64 //
	fromOffset   float64 //
	toOffset     float64 //
}
