package main

import "fmt"

type Equipment struct {
	ID    int    // database unique id
	Name  string // name of equipment
	Owned bool   // if the equipment is owned

}

func (e Equipment) String() string {

	if e.Owned {

		return fmt.Sprintf("Equipment: %s is owned", e.Name)
	} else {

		return fmt.Sprintf("Equipment: %s is not owned", e.Name)

	}
}

// ByName implements sort.Interface for []Equipment
// based on the Name field
// https://golang.org/pkg/sort/#pkg-overview
type ByName []Equipment

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j] }

// ByID implements sort.Interface for []Equipment
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByID []Equipment

func (d ByID) Len() int           { return len(d) }
func (d ByID) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByID) Less(i, j int) bool { return d[i].ID < d[j].ID }
