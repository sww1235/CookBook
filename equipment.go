package main

import "fmt"

type Equipment struct {
	ID    int    // database unique id
	Name  string // name of equipment
	Owned bool   // if the equipment is owned

}

func (e Equipment) String() string {

	if e.Owned {

		return fmt.Sprintf("Equipment: %s is owned")
	} else {

		return fmt.Sprintf("Equipment: %s is not owned")

	}
}
