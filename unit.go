package main

import "fmt"

type Unit struct {
	ID          int      // id of unit in database
	Name        string   // human readable name of unit
	Symbol      string   // recipe symbol
	Description string   // unit description
	IsCustom    bool     // is unit custom or standard
	UnitType    UnitType // type of unit

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

// ByIDUT implements sort.Interface for []UnitType
// based on the Name field
// https://golang.org/pkg/sort/#pkg-overview

type ByIDUT []UnitType

func (t ByIDUT) Len() int           { return len(t) }
func (t ByIDUT) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByIDUT) Less(i, j int) bool { return t[i].ID < t[j].ID }

// ByNameU implements sort.Interface for []Unit
// based on the Name field
// https://golang.org/pkg/sort/#pkg-overview
type ByNameU []Unit

func (n ByNameU) Len() int           { return len(n) }
func (n ByNameU) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByNameU) Less(i, j int) bool { return n[i].Name < n[j].Name }

// ByIDU implements sort.Interface for []Unit
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByIDU []Unit

func (d ByIDU) Len() int           { return len(d) }
func (d ByIDU) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByIDU) Less(i, j int) bool { return d[i].ID < d[j].ID }

// ByIDC implements sort.Interface for []Conversion
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByIDC []Conversion

func (d ByIDC) Len() int           { return len(d) }
func (d ByIDC) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByIDC) Less(i, j int) bool { return d[i].ID < d[j].ID }
