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
