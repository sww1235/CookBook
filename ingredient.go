package recipeDatabase

//The Ingredient struct stores data for a particular Ingredient
//used in a recipe
type Ingredient struct {
	Name               string
	UPC                string
	QuantityNeeded     int
	IngredientUnit     ingredientUnit
	InDatabase         bool
	QuantityInDatabase int
}

//Type ingredientUnit can be any value but only certain values are
//able to be converted at present. Values that are unable to be converted
//will not be able to be selected.
type ingredientUnit string
