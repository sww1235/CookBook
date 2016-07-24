package recipeDatabase

import "time"

//TODO: look at adding equipment/pots/pans to Recipe struct

//A Recipe struct is the internal representation of a recipe JSON file
type Recipe struct {
	Name              string
	FilePath          string
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
