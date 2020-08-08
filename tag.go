package main

type Tag struct {
	ID   int    // id of tag in database
	Name string // string representation of tag

}

func (t Tag) String() string {
	return t.Name

}
