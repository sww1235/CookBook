package recipeDatabase

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

//ReadRecipe creates a recipe struct by prompting the user for input
func ReadRecipe() (Recipe, error) {

	var tempRecipe Recipe
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter recipe Name: ")
	tempString, err := reader.ReadString('\n')
	if err != nil {
		return tempRecipe, err
	}
	tempRecipe.Name = tempString

	// read in quantity made and units
	fmt.Println("Enter in Quantity made and units in the following format: \"000 units\"")
	fmt.Println("Example entry: \"100 cookies\"")

	tempString, err = reader.ReadString('\n')
	if err != nil {
		return tempRecipe, err
	}
	qtyMadeReturn := strings.Split(tempString, " ")
	tempRecipe.QuantityMade, err = strconv.Atoi(qtyMadeReturn[0])
	if err != nil {
		return tempRecipe, err
	}
	tempRecipe.QuantityMadeUnits = qtyMadeReturn[1]

	// read in ingredients

	// read in steps

	return tempRecipe, nil
}
