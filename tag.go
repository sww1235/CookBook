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

func (t ByNameT) Len() int           { return len(t) }
func (t ByNameT) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByNameT) Less(i, j int) bool { return t[i].Name < t[j].Name }

// ByID implements sort.Interface for []Tag
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByIDT []Tag

func (t ByIDT) Len() int           { return len(t) }
func (t ByIDT) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByIDT) Less(i, j int) bool { return t[i].ID < t[j].ID }
