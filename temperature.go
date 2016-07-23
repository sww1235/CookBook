package recipeDatabase

//TempUnit is a single rune that represents the temperature scale associated
//with a particular temperature. The accepted options are F, C, K, R
//for Fehrenheit, Celsius, Kelvin, Rankine
type TempUnit rune

type temperature struct {
	Value float64
	Unit  TempUnit
}
