package main

import "fmt"

type Temperature struct {
	Value float64
	Unit  Unit
}

func (t Temperature) String() string {
	return fmt.Sprintf("%Gº %v", t.Value, t.Unit)
}
