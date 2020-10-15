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
type ByNameE []Equipment

func (n ByNameE) Len() int           { return len(n) }
func (n ByNameE) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByNameE) Less(i, j int) bool { return n[i].Name < n[j].Name }

// ByID implements sort.Interface for []Equipment
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByIDE []Equipment

func (d ByIDE) Len() int           { return len(d) }
func (d ByIDE) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByIDE) Less(i, j int) bool { return d[i].ID < d[j].ID }
