package main

import "fmt"

type Unit struct {
	ID            int        // id of unit in database
	Name          string     // human readable name of unit
	Symbol        string     // recipe symbol
	Description   string     // unit description
	IsCustom      bool       // is unit custom or standard
	RefIngredient Ingredient //Referenced ingredient for ingredient specific unit
	UnitType      UnitType   // type of unit

}

func (u Unit) String() string {
	stringString := ""
	stringString += fmt.Sprintf("Unit Name: %s\nUnit Symbol: %s\n", u.Name, u.Symbol)
	return stringString

}

// Type Conversion is a storage method for conversions between different units
// The conversions use the following formula.
// toUnitValue = ((fromUnitValue + fromOffset) * multiplicand / denominator) + toOffset
type Conversion struct {
	ID           int     // id of conversion factor in database
	FromUnit     int     // db id of unit to convert from
	ToUnit       int     // db id of unit to convert to
	Multiplicand float64 //
	Denominator  float64 //
	FromOffset   float64 //
	ToOffset     float64 //
}

// Type UnitType represents the 7 basic unit types as listed:
//
// time, length, mass, current, temperature, quantity, lum_intensity
type UnitType struct {
	ID   int
	Name string
}
