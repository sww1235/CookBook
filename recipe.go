package main

import (
	"fmt"
	"time"
)

//A Recipe struct is the internal representation of a recipe from a database
type Recipe struct {
	ID               int          // id of recipe in database, -1 if recipe doesn't exist in db
	Name             string       // name of recipe
	Description      string       // description of recipe
	Comments         string       // recipe comments
	Source           string       // source of recipe
	Author           string       // author of recipe
	Ingredients      []Ingredient // ingredients of recipe
	QuantityMade     float64      // how much of unit recipe makes
	QuantityMadeUnit Unit         // unit of recipe
	Steps            []Step       // steps of recipe
	EquipmentNeeded  []Equipment  // equipment needed to make recipe
	Tags             []Tag        // recipe tags
	Version          int          // version of recipe
	InitialVersion   int          // id of initial version of recipe
}

func (r Recipe) String() string {
	stringString := ""
	stringString += fmt.Sprintf("%s \n\n ", r.Name)
	if r.QuantityMade > 1 {
		stringString += fmt.Sprintf("Makes %G %s's\n", r.QuantityMade, r.QuantityMadeUnit)
	} else if r.QuantityMade == 1 {
		stringString += fmt.Sprintf("Makes %G %s\n", r.QuantityMade, r.QuantityMadeUnit)
	} else {
		stringString += "Makes nothing, good job cookie\n"
	}

	var prepTime, cookTime, waitTime, otherTime time.Duration = 0, 0, 0, 0

	for _, step := range r.Steps {
		switch step.StepType {
		case Prep:
			prepTime += step.TimeNeeded
		case Cook:
			cookTime += step.TimeNeeded
		case Wait:
			waitTime += step.TimeNeeded
		case Other:
			otherTime += step.TimeNeeded

		}

	}

	totalTime := prepTime + cookTime + waitTime + otherTime

	stringString += fmt.Sprintf("Takes: %d of total prep time\n", prepTime)
	stringString += fmt.Sprintf("Takes: %d of total cook time\n", cookTime)
	stringString += fmt.Sprintf("Takes: %d of total wait time\n", waitTime)
	stringString += fmt.Sprintf("Takes: %d of total other time\n", otherTime)
	stringString += fmt.Sprintf("Takes: %d of total time\n", totalTime)

	stringString += "Ingredients: \n"
	for _, Ingredient := range r.Ingredients {
		stringString += "\t" + Ingredient.String() //\n terminated
	}
	stringString += "Instructions: \n"
	for i, step := range r.Steps {
		stringString += fmt.Sprintf("%d) %s", i, step.String())
	}
	stringString += "Tags: \n"
	for _, tag := range r.Tags {
		stringString += tag.String()
	}

	stringString += "\n\n"
	return stringString
}

// ByName implements sort.Interface for []Recipe
// based on the Name field
// https://golang.org/pkg/sort/#pkg-overview
type ByNameR []Recipe

func (n ByNameR) Len() int           { return len(n) }
func (n ByNameR) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByNameR) Less(i, j int) bool { return n[i].Name < n[j] }

// ByID implements sort.Interface for []Recipe
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByIDR []Recipe

func (d ByIDR) Len() int           { return len(d) }
func (d ByIDR) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByIDR) Less(i, j int) bool { return d[i].ID < d[j].ID }
