package main

import "fmt"

type temperature struct {
	Value float64
	Unit  Unit
}

func (t temperature) String() string {
	return fmt.Sprintf("%Gº %v", t.Value, t.Unit)
}
