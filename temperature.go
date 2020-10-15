package main

import "fmt"

type Temperature struct {
	Value float64
	Unit  Unit
}

func (t Temperature) String() string {
	return fmt.Sprintf("%Gº %v", t.Value, t.Unit)
}

// ByName implements sort.Interface for []Temperature
// based on the Name field
// https://golang.org/pkg/sort/#pkg-overview
type ByValueT []Temperature

func (t ByValueT) Len() int           { return len(t) }
func (t ByValueT) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByValueT) Less(i, j int) bool { return t[i].Value < t[j].Value }
