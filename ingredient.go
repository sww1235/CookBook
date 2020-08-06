package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

//The Ingredient struct stores data for a particular Ingredient
//used in a recipe
type Ingredient struct {
	Name               string
	UPC                string
	QuantityNeeded     float64
	IngredientUnit     ingredientUnit
	InDatabase         bool
	QuantityInDatabase int
	Conversions        []conversion
}

//Type ingredientUnit can be any value but only certain values are
//able to be converted at present. Values that are unable to be converted
//will not be able to be selected.
type ingredientUnit string

func (i Ingredient) String() string {
	stringString := fmt.Sprintf("%s: %G %s(s)\n", i.Name, i.QuantityNeeded, i.IngredientUnit)

	return stringString
}

//ConvertString acts like String() but allows for conversion between units
func (i Ingredient) ConvertString(toUnit string) string {
	return "fixme" //TODO: implement
}

//AddConversion adds a conversion factor to an ingredient
func (i *Ingredient) AddConversion(fromUnit string, toUnit string, factor float64) {

	i.Conversions = append(i.Conversions, conversion{fromUnit, toUnit, factor})
}

//type conversion is a storage method for conversions between different units
//ConversionFactor is the number to multiply by to convert from FromUnit to ToUnit
type conversion struct {
	FromUnit         string
	ToUnit           string
	ConversionFactor float64
}

// ReadIngredient creates an ingredient struct by prompting user for input
func ReadIngredient() (Ingredient, error) {
	var tempIngredient Ingredient

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter ingredient Name: ")
	tempString, err := reader.ReadString('\n')
	if err != nil {
		return tempIngredient, err
	}
	tempIngredient.Name = tempString

	fmt.Print("Enter ingredient UPC: ")
	tempString, err = reader.ReadString('\n')
	if err != nil {
		return tempIngredient, err
	}
	tempIngredient.UPC = tempString

	fmt.Print("Enter ingredient Quantity: ")
	tempString, err = reader.ReadString('\n')
	if err != nil {
		return tempIngredient, err
	}
	tempQty, err := strconv.ParseFloat(tempString, 64)

	tempIngredient.QuantityNeeded = tempQty

	return tempIngredient, nil

}
