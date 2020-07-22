package recipeDatabase

import (
	"fmt"
	"time"
)

//TODO: look at adding equipment/pots/pans to Recipe struct

//A Recipe struct is the internal representation of a recipe JSON file
type Recipe struct {
	Name              string
	Ingredients       []Ingredient
	QuantityMade      int
	QuantityMadeUnits string
	Steps             []step
	TotalPrepTime     time.Duration
	TotalCookTime     time.Duration
	TotalWaitTime     time.Duration
	TotalOtherTime    time.Duration
	TotalTime         time.Duration
	Tags              []string
}

//StepType has 4 recognized values, prep, cook, wait
//and then anything else as other. This selects what duration in a Recipe
//TimeNeeded is added to.
type StepType string

type step struct {
	TimeNeeded   time.Duration
	StepType     StepType
	Temperature  temperature
	Instructions string
}

func (s step) String() string {
	stringString := ""
	stringString += fmt.Sprintf("%s: Needs %d\nCook at %s\n", s.StepType, s.TimeNeeded, s.Temperature.String())
	stringString += s.Instructions + "\n"

	return stringString
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
	stringString += fmt.Sprintf("Takes: %d of total prep time\n", r.TotalPrepTime)
	stringString += fmt.Sprintf("Takes: %d of total cook time\n", r.TotalCookTime)
	stringString += fmt.Sprintf("Takes: %d of total wait time\n", r.TotalWaitTime)
	stringString += fmt.Sprintf("Takes: %d of total other time\n", r.TotalOtherTime)
	stringString += fmt.Sprintf("Takes: %d of total time\n", r.TotalTime)

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
