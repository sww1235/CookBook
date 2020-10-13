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
type ByNameT []Tag

func (n ByNameT) Len() int           { return len(n) }
func (n ByNameT) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByNameT) Less(i, j int) bool { return n[i].Name < n[j] }

// ByID implements sort.Interface for []Tag
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByIDT []Tag

func (d ByIDT) Len() int           { return len(d) }
func (d ByIDT) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByIDT) Less(i, j int) bool { return d[i].ID < d[j].ID }
