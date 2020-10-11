package main

type Tag struct {
	ID          int    // id of tag in database
	Name        string // string representation of tag
	Description string // tag description

}

func (t Tag) String() string {
	return t.Name + ": " + t.Description

}
