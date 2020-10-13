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
	ID                 int             // id in database
	Name               string          // ingredient common name
	QuantityUsed       float64         // quantity of ingredient in recipe
	QuantityUnit       Unit            // unit of ingredient quantity in recipe
	InventoryReference []InventoryItem // inventory items that creates ingredient
}

func (i Ingredient) String() string {

	return fmt.Sprintf("%s: %G %s(s)\n", i.Name, i.QuantityUsed, i.QuantityUnit)

}

// ByName implements sort.Interface for []Ingredient
// based on the Name field
// https://golang.org/pkg/sort/#pkg-overview
type ByNameIn []Ingredient

func (n ByNameIn) Len() int           { return len(n) }
func (n ByNameIn) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByNameIn) Less(i, j int) bool { return n[i].Name < n[j] }

// ByID implements sort.Interface for []Ingredient
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByIDIn []Ingredient

func (d ByIDIn) Len() int           { return len(d) }
func (d ByIDIn) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByIDIn) Less(i, j int) bool { return d[i].ID < d[j].ID }

//// ConvertString acts like String() but allows for conversion between units
//func (i Ingredient) ConvertString(toUnit string) string {
//	return "fixme" //TODO: implement
//}
//
//// Convert converts Ingredient quantityUsed from default QuantityUnit to toUnit
//func (i *Ingredient) Convert(toUnit Unit) (float64, error) {
//	var uC conversion
//	for _, conversion := range i.Conversions {
//
//		if conversion.FromUnit == i.QuantityUnit.ID && conversion.ToUnit == toUnit.ID {
//			uC = conversion
//		}
//	}
//	// check if conversion factor wasn't found, and throw an error
//	// compare to empty conversion struct
//	if (uC == conversion{}) {
//		return 0, fmt.Errorf("Conversion factor for unit pair %s:%s not found in list of conversion factors for ingredient %s",
//			i.QuantityUnit.Name, toUnit.Name, i.Name)
//	}
//
//	return ((i.QuantityUsed + uC.fromOffset) * uC.multiplicand / uC.denominator) + uC.toOffset, nil
//
//}

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
