package main

import "fmt"

//TempUnit is a single rune that represents the temperature scale associated
//with a particular temperature. The accepted options are F, C, K, R
//for Fehrenheit, Celsius, Kelvin, Rankine
type TempUnit rune

type temperature struct {
	Value float64
	Unit  TempUnit
}

func (t temperature) String() string {
	return fmt.Sprintf("%Gº %v", t.Value, t.Unit)
}

func (t temperature) Convert(dest rune) float64 {
	switch dest {
	case 'F':
		switch t.Unit {
		case 'F':
			return t.Value
		case 'C':
			return (t.Value * 9.0 / 5.0) + 32
		case 'K':
			return (t.Value * 9.0 / 5.0) - 459.67
		case 'R':
			return t.Value - 459.67
		}
	case 'C':
		switch t.Unit {
		case 'F':
			return (t.Value - 32) * (5.0 / 9.0)
		case 'C':
			return t.Value
		case 'K':
			return t.Value - 273.15
		case 'R':
			return (t.Value - 491.67) * (5.0 / 9.0)
		}
	case 'K':
		switch t.Unit {
		case 'F':
			return (t.Value + 459.67) * (5.0 / 9.0)
		case 'C':
			return t.Value + 273.15
		case 'K':
			return t.Value
		case 'R':
			return t.Value * (5.0 / 9.0)
		}
	case 'R':
		switch t.Unit {
		case 'F':
			return (t.Value + 459.67)
		case 'C':
			return (t.Value + 273.15) * (9.0 / 5.0)
		case 'K':
			return t.Value * (9.0 / 5.0)
		case 'R':
			return t.Value
		}
	default:
		fmt.Println("unrecognized temp unit")
		return -1
	}
	fmt.Println("unrecognized temp unit")
	return -1
}
