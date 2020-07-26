package recipeDatabase

type Unit struct {
	ID     int    // id of unit in database
	Name   string // human readable name of unit
	Symbol string // recipe symbol

}

func (u Unit) String() string {
	stringString := ""
	stringString += fmt.Sprintf("Unit Name: %s\nUnit Symbol: %s\n", u.Name, u.Symbol)

}
