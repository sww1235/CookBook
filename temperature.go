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
type ByName []Temperature

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j] }

// ByID implements sort.Interface for []Temperature
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByID []Temperature

func (d ByID) Len() int           { return len(d) }
func (d ByID) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByID) Less(i, j int) bool { return d[i].ID < d[j].ID }
