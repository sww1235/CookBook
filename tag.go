package main

type Tag struct {
	ID          int    // id of tag in database
	Name        string // string representation of tag
	Description string // tag description

}

func (t Tag) String() string {
	return t.Name + ": " + t.Description

}

// ByName implements sort.Interface for []Tag
// based on the Name field
// https://golang.org/pkg/sort/#pkg-overview
type ByName []Tag

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j] }

// ByID implements sort.Interface for []Tag
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByID []Tag

func (d ByID) Len() int           { return len(d) }
func (d ByID) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByID) Less(i, j int) bool { return d[i].ID < d[j].ID }
