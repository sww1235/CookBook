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

// ByName implements sort.Interface for []Unit
// based on the Name field
// https://golang.org/pkg/sort/#pkg-overview
type ByName []Unit

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j] }

// ByID implements sort.Interface for []Unit
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByID []Unit

func (d ByID) Len() int           { return len(d) }
func (d ByID) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByID) Less(i, j int) bool { return d[i].ID < d[j].ID }

// ByID implements sort.Interface for []Conversion
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByID []Conversion

func (d ByID) Len() int           { return len(d) }
func (d ByID) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByID) Less(i, j int) bool { return d[i].ID < d[j].ID }
