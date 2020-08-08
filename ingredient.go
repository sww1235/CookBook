package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// The Ingredient struct stores data for a particular Ingredient
// used in a recipe
type Ingredient struct {
	ID                 int           // id in database
	Name               string        // ingredient common name
	QuantityUsed       float64       // quantity of ingredient in recipe
	QuantityUnit       Unit          // unit of ingredient quantity in recipe
	InventoryReference InventoryItem // inventory item that creates ingredient
	Conversions        []conversion  // in memory store of conversions associated with ingredient
}

func (i Ingredient) String() string {

	return fmt.Sprintf("%s: %G %s(s)\n", i.Name, i.QuantityUsed, i.QuantityUnit)

}

// ConvertString acts like String() but allows for conversion between units
func (i Ingredient) ConvertString(toUnit string) string {
	return "fixme" //TODO: implement
}

// AddConversion adds a conversion factor to an ingredient
func (i *Ingredient) AddConversion(id int, fromUnit Unit, toUnit Unit, multiplicand float64,
	denominator float64, fromOffset float64, toOffset float64) {

	i.Conversions = append(i.Conversions, conversion{id, fromUnit.ID, toUnit.ID,
		multiplicand, denominator, fromOffset, toOffset})
}

// Convert converts Ingredient quantityUsed from default QuantityUnit to toUnit
func (i *Ingredient) Convert(toUnit Unit) (float64, error) {
	var uC conversion
	for _, conversion := range i.Conversions {

		if conversion.FromUnit == i.QuantityUnit.ID && conversion.ToUnit == toUnit.ID {
			uC = conversion
		}
	}
	// check if conversion factor wasn't found, and throw an error
	// compare to empty conversion struct
	if (uC == conversion{}) {
		return 0, fmt.Errorf("Conversion factor for unit pair %s:%s not found in list of conversion factors for ingredient %s",
			i.QuantityUnit.Name, toUnit.Name, i.Name)
	}

	return ((i.QuantityUsed + uC.fromOffset) * uC.multiplicand / uC.denominator) + uC.toOffset, nil

}

// LoadConversionFactors loads conversion factors for the referenced unit
// into the Ingredients Conversions list from the conversionFactors table
// in the database
func (i *Ingredient) LoadConversionFactors() error {

	return nil
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

	fmt.Print("Enter ingredient Quantity: ")
	tempString, err = reader.ReadString('\n')
	if err != nil {
		return tempIngredient, err
	}
	tempQty, err := strconv.ParseFloat(tempString, 64)

	tempIngredient.QuantityUsed = tempQty

	return tempIngredient, nil

}
