package recipeDatabase

import "time"

//A Recipe struct is the internal representation of a recipe JSON file
type Recipe struct {
	Ingredients       []Ingredient
	QuantityMade      int
	QuantityMadeUnits string
	Steps             []step
}

type step struct {
	TimeNeeded  time.Duration
	Temperature temperature
}
