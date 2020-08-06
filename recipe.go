package main

import (
	"fmt"
	"time"
)

//TODO: look at adding equipment/pots/pans to Recipe struct

//A Recipe struct is the internal representation of a recipe from a database
type Recipe struct {
	ID                int          // id of recipe in database
	Name              string       // name of recipe
	Description       string       // description of recipe
	Comments          string       // recipe comments
	Source            string       // source of recipe
	Author            string       // author of recipe
	Ingredients       []Ingredient // ingredients of recipe
	QuantityMade      int          // how much of unit recipe makes
	QuantityMadeUnits Unit         // unit of recipe
	Steps             []Step       // steps of recipe
	EquipmentNeeded   []Equipment  // equipment needed to make recipe
	Tags              []string     // recipe tags
}

func (r Recipe) String() string {
	stringString := ""
	stringString += fmt.Sprintf("%s \n\n ", r.Name)
	if r.QuantityMade > 1 {
		stringString += fmt.Sprintf("Makes %d %s's\n", r.QuantityMade, r.QuantityMadeUnits)
	} else if r.QuantityMade == 1 {
		stringString += fmt.Sprintf("Makes %d %s\n", r.QuantityMade, r.QuantityMadeUnits)
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
		stringString += tag
	}

	stringString += "\n\n"
	return stringString
}
