package main

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//initDB opens a connection to a sqlite database that stores recipes
func initDB(databasePath string) (*sql.DB, error) {

	// TODO: allow for other database backends, and allow for tables to be in separate backends

	// try to create database directory, then check if database exists before opening.
	// race condition exists on check if file exists, but I don't think this application
	// will run into it.

	mkErr := os.MkdirAll(path.Dir(databasePath), 0744)
	if mkErr != nil {
		return nil, fmt.Errorf("Unable to use %s as recipe directory, Err: %w",
			path.Dir(databasePath), mkErr)
	}

	_, err := os.Stat(databasePath)
	if os.IsNotExist(err) {
		infoLogger.Printf("database doesn't exist, creating now at path %s", databasePath)
	}

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("Could not open recipe database. %w", err)
	}

	// now check if correct tables exist, and create them if they do not.

	needInit := true
	missingTable := false
	requiredTables := map[string]bool{
		"recipes":              false,
		"ingredients":          false,
		"ingredient_inventory": false,
		"ingredient_recipe":    false,
		"steps":                false,
		//"stepType":             false,
		"step_recipe":      false,
		"inventory":        false,
		"units":            false,
		"tags":             false,
		"tag_recipe":       false,
		"lastMade":         false,
		"equipment":        false,
		"equipment_recipe": false,
		"unitType":         false,
		"unitConversions":  false,
	}

	for table := range requiredTables {
		sqlStatement := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"

		var rowCount int
		err := db.QueryRow(sqlStatement, table).Scan(&rowCount)
		if err != nil {
			return db, fmt.Errorf("could not check if table exists,%w", err)
		}

		debugLogger.Println(rowCount)

		if rowCount == 0 {
			requiredTables[table] = true

		}
	}

	for table, missing := range requiredTables {
		// any false values will cancel out the true initial value
		// for needInit. needInit will only be true if all tables are missing
		needInit = missing && needInit
		// missingTable will be true if any of the tables are missing
		missingTable = missing || missingTable
		if missing {
			infoLogger.Printf("Table %s missing.", table)
		}
	}

	// only needs to happen if some tables are missing, not all.
	// if all tables are missing, then needsInit is true
	if missingTable && !needInit {
		return db, fmt.Errorf("Existing database missing at least one critical table. See log messages above." +
			"Either delete database so it can be recreated automatically, or manually create the missing tables")
	}

	if needInit {
		// need to create tables
		// using map to allow for easier iteration
		createQueries := make(map[string]string)

		// id column will automap to rowID autocreated by sqlite with "INTEGER PRIMARY KEY" spec.

		createQueries["equipment"] = "CREATE TABLE equipment(id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, isOwned NUM)"

		createQueries["unitType"] = "CREATE TABLE unitType(id INTEGER NOT NULL PRIMARY KEY, name TEXT)"

		createQueries["tags"] = "CREATE TABLE tags(id INTEGER NOT NULL PRIMARY KEY, name TEXT NOT NULL, " +
			"description TEXT)"

		//createQueries["stepType"] = "CREATE TABLE stepType (id INTEGER NOT NULL PRIMARY KEY, " +
		//	"name TEXT)"

		createQueries["units"] = "CREATE TABLE units(id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, description TEXT, symbol TEXT, isCustom INTEGER, " +
			"refIngredient INTEGER, unitType INTEGER, " +
			"FOREIGN KEY (refIngredient) REFERENCES ingredients(id), " +
			"FOREIGN KEY (unitType) REFERENCES unitType(id))"

		createQueries["recipes"] = "CREATE TABLE recipes (id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, description TEXT, comments TEXT, source TEXT, author TEXT, " +
			"quantity NUM, quantityUnits INTEGER, initialVersion INTEGER, version INTEGER, " +
			"FOREIGN KEY(quantityUnits) REFERENCES units(id), " +
			"FOREIGN KEY(initialVersion) REFERENCES recipes(id))"

		createQueries["inventory"] = "CREATE TABLE inventory (id INTEGER NOT NULL PRIMARY KEY, " +
			"EAN TEXT UNIQUE, name TEXT, description TEXT, quantity NUM, " +
			"QuantityUnits INTEGER, FOREIGN KEY(QuantityUnits) REFERENCES units(id))"

		createQueries["ingredients"] = "CREATE TABLE ingredients (id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, quantity NUM, quantityUnits INTEGER, " +
			"FOREIGN KEY(quantityUnits) REFERENCES units(id))"

		createQueries["ingredient_inventory"] = "CREATE TABLE ingredient_inventory( " +
			"ingredientID INTEGER NOT NULL, inventoryID INTEGER NOT NULL, " +
			"FOREIGN KEY(ingredientID) REFERENCES ingredients(id), " +
			"FOREIGN KEY(inventoryID) REFERENCES inventory(id), " +
			"PRIMARY KEY(ingredientID, inventoryID))"

		createQueries["ingredient_recipe"] = "CREATE TABLE ingredient_recipe( " +
			"ingredientID INTEGER NOT NULL, recipeID INTEGER NOT NULL, " +
			"FOREIGN KEY(ingredientID) REFERENCES ingredients(id), " +
			"FOREIGN KEY(recipeID) REFERENCES recipes(id), " +
			"PRIMARY KEY(ingredientID, recipeID))"

		createQueries["steps"] = "CREATE TABLE steps( id INTEGER NOT NULL PRIMARY KEY, order INTEGER " +
			"instructions TEXT, time NUM, stepType INTEGER, temperature NUM, tempUnit INTEGER, " +
			"FOREIGN KEY(stepTypeID) REFERENCES stepType(id), " +
			"FOREIGN KEY(tempUnits) REFERENCES units(id))"

		createQueries["step_recipe"] = "CREATE TABLE step_recipe( stepID INTEGER NOT NULL, " +
			"recipeID INTEGER NOT NULL, " +
			"FOREIGN KEY(stepID) REFERENCES steps(id), " +
			"FOREIGN KEY(recipeID) REFERENCES recipes(id), " +
			"PRIMARY KEY(stepID, recipeID))"

		createQueries["tag_recipe"] = "CREATE TABLE tag_recipe( tagID INTEGER NOT NULL, recipeID INTEGER NOT NULL, " +
			"FOREIGN KEY(tagID) REFERENCES tags(id), " +
			"FOREIGN KEY(recipeID) REFERENCES recipes(id), " +
			"PRIMARY KEY(tagID, recipeID))"

		createQueries["equipment_recipe"] = "CREATE TABLE equipment_recipe( equipmentID INTEGER NOT NULL, " +
			"recipeID INTEGER NOT NULL, " +
			"FOREIGN KEY(equipmentID) REFERENCES equipment(id), " +
			"FOREIGN KEY(recipeID) REFERENCES recipes(id), " +
			"PRIMARY KEY(equipmentID, recipeID))"

		createQueries["lastMade"] = "CREATE TABLE lastMade(id INTEGER NOT NULL PRIMARY KEY, " +
			"recipe INTEGER NOT NULL, dateMade TEXT, notes TEXT, " +
			"FOREIGN KEY(recipe) REFERENCES recipes(id))"

		createQueries["unitConversions"] = "CREATE TABLE unitConversions( id INTEGER NOT NULL PRIMARY KEY, " +
			"fromUnit INTEGER, toUnit INTEGER, multiplicand NUM, denominator NUM, fromOFfset NUM, toOffset NUM, " +
			"FOREIGN KEY(fromUnit) REFERENCES units(id), " +
			"FOREIGN KEY(toUnit) REFERENCES units(id))"

		// since not all tables exist, for now drop all tables, then recreate them

		// now create all the tables
		for table, query := range createQueries {

			_, err := db.Exec(query)
			if err != nil {
				return db, fmt.Errorf("Failed to create table: %s due to error: %w", table, err)
			}
		}

	}
	return db, nil

}

// InsertRecipe takes a Recipe struct of a recipe that doesn't exist in the database
// and inserts its contents into the appropriate tables in the database
//
// When using this function, the Recipe struct should have all fields filled out, including
// sub structs such as ingredients and steps.
//
// Returns id of new recipe in database, and any associated errors
func InsertRecipe(db *sql.DB, r Recipe) (int, error) {

	// TODO: switch from db.Exec to transactions that are chained through all insert/update functions
	// need to insert recipe first, so steps can be linked to appropriate ID
	insertSQL := "INSERT INTO recipes (name, description, comments, source, author, quantity, " +
		"quantityUnits, initialVersion, version) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"

	// if recipe already exists in database
	if r.ID != -1 {
		return -1, fmt.Errorf("Recipe %s already exists in database, cannot insert duplicate", r)
	}

	result, err := db.Exec(insertSQL, r.Name, r.Description, r.Comments, r.Source, r.Author,
		r.QuantityMade, r.QuantityMadeUnit, r.InitialVersion, r.Version)

	if err != nil {
		return -1, err
	}

	recipeID, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	for _, ing := range r.Ingredients {
		// ingredient_recipe mapping is done in InsertIngredient
		_, err := InsertIngredient(db, ing, r, ing.InventoryReference)

		if err != nil {
			return -1, err
		}
	}

	//TODO: add steps, equipment, and tags

	return int(recipeID), err
}

// UpdateRecipe takes a Recipe struct that exists in the database, and updates all fields
// to be the same as in the passed in Recipe struct.
//
// When using this function, the Recipe struct should have all fields filled out, including
// sub structs such as ingredients and steps.
//
// Returns id of updated recipe in database, and any associated errors
func UpdateRecipe(db *sql.DB, r Recipe) (int, error) {

	// don't need to update initialVersion
	updateSQL := "UPDATE recipes SET name = ?, description = ?, comments = ?, source = ?, author = ?, " +
		"quantity = ?, quantityUnits = ?, version = ? " +
		"WHERE id = ?"

	// if recipe does not exist in database
	if r.ID == -1 {
		return -1, fmt.Errorf("Cannot update recipe %s, it does not exist", r)
	}

	_, err := db.Exec(updateSQL, r.Name, r.Description, r.Comments, r.Source, r.Author,
		r.QuantityMade, r.QuantityMadeUnit, r.Version, r.ID)

	return r.ID, err

}

// InsertIngredient takes an Ingredient struct of an ingredient that doesn't exist in the database
// and inserts its contents into the appropriate tables in the database.
//
// When using this function, the Ingredient struct should have all fields filled out.
//
// Returns id of inserted ingredient in database, and any associated errors
func InsertIngredient(db *sql.DB, ing Ingredient, r Recipe, invUnitID int) (int, error) {

	insertSQL := "INSERT INTO ingredients (name, quantity, quantityUnits) " +
		"VALUES(?, ?, ?)"

	mapRecSQL := "INSERT INTO ingredient_recipe (recipeID, ingredientID) VALUES (?, ?)"

	mapInvSQL := "INSERT INTO ingredient_inventory (ingredientID, inventoryID) VALUES (?, ?)"

	// if ingredient already exists in database
	if ing.ID != -1 {
		return -1, fmt.Errorf("Ingredient %s already exists in database", ing)
	}

	result, err := db.Exec(insertSQL, ing.Name, ing.QuantityUsed, ing.QuantityUnit.ID)

	if err != nil {
		return -1, err
	}

	ingredientID, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	_, err = db.Exec(mapRecSQL, r.ID, ingredientID)

	if err != nil {
		return -1, err
	}

	_, err = db.Exec(mapInvSQL, ingredientID, invUnitID)

	if err != nil {
		return -1, err
	}

	return int(ingredientID), err

}

// UpdateIngredient takes an Ingredient struct of an ingredient that exists in the database
// and updates the database with the values in the struct.
//
// When using this function, the Ingredient struct should have all fields filled out.
//
// Returns id of updated ingredient in database, and any associated errors
func UpdateIngredient(db *sql.DB, ing Ingredient, invUnitID int) (int, error) {

	updateIngredientSQL := "UPDATE ingredients SET name = ?, quantity = ?, quantityUnits = ? " +
		"WHERE id = ?"

	// if ingredient doesn't exist in database
	if ing.ID == -1 {
		return -1, fmt.Errorf("Ingredient %s doesn't exist in database, can't update it", ing)

	}

	_, err := db.Exec(updateIngredientSQL, ing.Name, ing.QuantityUsed, ing.QuantityUnit.ID, ing.ID)

	//TODO: need to delete and remap ingredient_inventory

	return ing.ID, err
}

// InsertStep takes an Step struct of a step that does not exist in the database
// and inserts its contents into the appropriate tables in the database
//
// When using this function, the Step struct must have all fields filled out except for temperature and tempUnit.
//
// Returns id of inserted step in database, and any associated errors.
func InsertStep(db *sql.DB, stp Step, r Recipe) (int, error) {

	insertSQL := "INSERT INTO steps (order, instructions, time, stepType, temperature, tempUnit) " +
		"VALUES(?, ?, ?, ?, ?, ?)"

	mapSQL := "INSERT INTO step_recipe (recipeID, stepID) VALUES (?, ?)"

	// if ingredient already exists in database
	if stp.ID != -1 {
		return -1, fmt.Errorf("Step %s already exists in database", stp)
	}

	result, err := db.Exec(insertSQL, stp.Order, stp.Instructions, stp.TimeNeeded, stp.StepType,
		stp.Temperature.Value, stp.Temperature.Unit.ID)

	if err != nil {
		return -1, err
	}

	stepID, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	_, err = db.Exec(mapSQL, r.ID, stepID)

	if err != nil {
		return -1, err
	}

	return int(stepID), err

}

// InsertStep takes an Step struct of a step that does not exist in the database
// and inserts its contents into the appropriate tables in the database
//
// When using this function, the Step struct must have all fields filled out except for temperature and tempUnit.
//
// Returns id of inserted step in database, and any associated errors.
func UpdateStep(db *sql.DB, stp Step) (int, error) {

	updateSQL := "UPDATE steps SET order = ?, instructions = ?, time = ?, stepType = ? " +
		"temperature = ?, tempUnit = ? WHERE id = ?"

	// if ingredient doesn't exist in database
	if stp.ID == -1 {
		return -1, fmt.Errorf("Step %s doesn't exist in database, can't update it", stp)

	}

	_, err := db.Exec(updateSQL, stp.Order, stp.Instructions, stp.TimeNeeded, stp.StepType,
		stp.Temperature.Value, stp.Temperature.Unit.ID, stp.ID)

	//TODO: need to delete and remap ingredient_inventory

	return stp.ID, err
}

// GetRecipes returns a slice of all recipes in database
func GetRecipes(db *sql.DB) ([]Recipe, error) {

	sqlString := "SELECT id, name, description, comments, source, author, quantity, " +
		"quantityMadeUnits, initialVersion, version FROM recipes"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	units, err := GetUnits(db)
	if err != nil {
		return nil, err
	}

	sort.Stable(ByIDU(units)) // sort units slice stably

	for rows.Next() {
		var (
			id             int
			name           string
			description    string
			comments       string
			source         string
			author         string
			quantity       float64
			qtyUnitsId     int //fk
			initialVersion int //fk
			version        int
		)

		err = rows.Scan(&id, &name, &description, &comments, &source, &author,
			&quantity, &qtyUnitsId, &initialVersion, &version)

		if err != nil {
			return nil, err
		}
		ingredients, err := GetIngredientsForRecipe(db, id) // returns slice of ingredients
		if err != nil {
			return nil, err
		}
		steps, err := GetStepsForRecipe(db, id) // returns slice of steps
		if err != nil {
			return nil, err
		}
		equipment, err := GetEquipmentForRecipe(db, id) // returns slice of equipment
		if err != nil {
			return nil, err
		}
		tags, err := GetTagsForRecipe(db, id) // returns slice of tags
		if err != nil {
			return nil, err
		}

		// https://golang.org/pkg/sort/#Search
		// first find index
		idIdx := sort.Search(len(units), func(i int) bool { return units[i].ID >= qtyUnitsId })
		// then perform sanity checks
		var tempUnit Unit
		if idIdx < len(units) && units[idIdx].ID == qtyUnitsId {
			tempUnit = units[idIdx]
		} else {
			return nil, fmt.Errorf("ID: %d not found in units slice %v", qtyUnitsId, units)
		}

		debugLogger.Printf("Recipe ID: %d, Recipe Name: %s", id, name)
		tempRecipe := Recipe{ID: id, Name: name, Description: description, Comments: comments,
			Source: source, Author: author, Ingredients: ingredients, QuantityMade: quantity,
			QuantityMadeUnit: tempUnit, Steps: steps,
			EquipmentNeeded: equipment, Tags: tags, Version: version, InitialVersion: initialVersion}
		recipes = append(recipes, tempRecipe)
	}

	return recipes, nil

}

// GetEquipmentForRecipe returns a map of Equipment structs indexed on
// the database table index for the specified recipe
func GetEquipmentForRecipe(db *sql.DB, recipeID int) ([]Equipment, error) {

	sqlString := "SELECT id, name, isOwned FROM equipment AS equip INNER JOIN equipment_recipe AS er" +
		"ON equip.id = er.equipmentID WHERE er.recipeID = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipment []Equipment

	for rows.Next() {
		var (
			id      int
			name    string
			isOwned bool
		)

		err = rows.Scan(&id, &name, &isOwned)
		if err != nil {
			return nil, err
		}

		debugLogger.Printf("Ingredient ID: %d, Ingredient Name: %s", id, name)
		tempEquipment := Equipment{ID: id, Name: name, Owned: isOwned}
		equipment = append(equipment, tempEquipment)
	}

	return equipment, nil
}

// GetIngredientsForRecipe returns a map of Ingredient structs indexed on
// the database table index for the specified recipe
func GetIngredientsForRecipe(db *sql.DB, recipeID int) ([]Ingredient, error) {

	sqlString := "SELECT id, name, quantity, quantityUnits FROM ingredients AS ing " +
		"INNER JOIN ingredient_recipe AS ir ON ing.id = ir.ingredientID WHERE ir.recipeID = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []Ingredient

	units, err := GetUnits(db)
	if err != nil {
		return nil, err
	}

	sort.Stable(ByIDU(units)) // sort units slice stably

	for rows.Next() {
		var (
			id       int
			name     string
			quantity float64
			unit     int //fk
		)

		err = rows.Scan(&id, &name, &quantity, &unit)
		if err != nil {
			return nil, err
		}
		inventory, err := GetInventoryForIngredient(db, id) // returns inventory struct
		if err != nil {
			return nil, err
		}
		// https://golang.org/pkg/sort/#Search
		// first find index
		idIdx := sort.Search(len(units), func(i int) bool { return units[i].ID >= unit })

		// then perform sanity checks
		var tempUnit Unit
		if idIdx < len(units) && units[idIdx].ID == unit {
			tempUnit = units[idIdx]
		} else {
			return nil, fmt.Errorf("ID: %d not found in units slice %v", unit, units)
		}

		debugLogger.Printf("Ingredient ID: %d, Ingredient Name: %s", id, name)
		tempIngredient := Ingredient{ID: id, Name: name, QuantityUsed: quantity,
			QuantityUnit: tempUnit, InventoryReference: inventory}
		ingredients = append(ingredients, tempIngredient)
	}

	return ingredients, nil
}

// GetStepsForRecipe returns a slice of Step structs
// for the specified recipe
func GetStepsForRecipe(db *sql.DB, recipeID int) ([]Step, error) {

	sqlString := "SELECT id, order, instructions, time, stepType, temperature, tempUnit " +
		"FROM steps INNER JOIN step_recipe AS isi " +
		"ON steps.id = is.stepID WHERE ir.recipeID = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []Step

	units, err := GetUnits(db)
	if err != nil {
		return nil, err
	}

	sort.Stable(ByIDU(units)) // sort units slice stably

	for rows.Next() {
		var (
			id           int
			order        int
			instructions string
			timeNeeded   int //seconds
			stepType     int //fk to stepType iota
			temperature  float64
			tempUnit     int //fk
		)

		err = rows.Scan(&id, &order, &instructions, &timeNeeded, &stepType,
			&temperature, &tempUnit)
		if err != nil {
			return nil, err
		}

		// https://golang.org/pkg/sort/#Search
		// first find index
		idIdx := sort.Search(len(units), func(i int) bool { return units[i].ID >= tempUnit })

		// then perform sanity checks
		var temperatureUnit Unit
		if idIdx < len(units) && units[idIdx].ID == tempUnit {
			temperatureUnit = units[idIdx]
		} else {
			return nil, fmt.Errorf("ID: %d not found in units slice %v", tempUnit, units)
		}

		duration, err := time.ParseDuration(strconv.Itoa(timeNeeded) + "s")
		if err != nil {
			return nil, err
		}
		temp := Temperature{Value: temperature, Unit: temperatureUnit}

		debugLogger.Printf("Step ID: %d", id)
		tempStep := Step{ID: id, Order: order, TimeNeeded: duration, StepType: StepType(stepType),
			Temperature: temp}
		steps = append(steps, tempStep)
	}

	return steps, nil
}

// GetIngredients returns a slice of all ingredients in the databse
func GetIngredients(db *sql.DB) ([]Ingredient, error) {

	sqlString := "SELECT id, name, quantity, quantityUnits FROM ingredients"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []Ingredient

	units, err := GetUnits(db)
	if err != nil {
		return nil, err
	}

	sort.Stable(ByIDU(units)) // sort units slice stably

	for rows.Next() {
		var (
			id       int
			name     string
			quantity float64
			unit     int //fk
		)

		err = rows.Scan(&id, &name, &quantity, &unit)
		if err != nil {
			return nil, err
		}
		inventory, err := GetInventoryForIngredient(db, id) // returns inventory struct
		if err != nil {
			return nil, err
		}

		// https://golang.org/pkg/sort/#Search
		// first find index
		idIdx := sort.Search(len(units), func(i int) bool { return units[i].ID >= unit })

		// then perform sanity checks
		var tempUnit Unit
		if idIdx < len(units) && units[idIdx].ID == unit {
			tempUnit = units[idIdx]
		} else {
			return nil, fmt.Errorf("ID: %d not found in units slice %v", unit, units)
		}

		debugLogger.Printf("Ingredient ID: %d, Ingredient Name: %s", id, name)
		tempIngredient := Ingredient{ID: id, Name: name, QuantityUsed: quantity,
			QuantityUnit: tempUnit, InventoryReference: inventory}
		ingredients = append(ingredients, tempIngredient)
	}

	return ingredients, nil
}

// GetConversions returns a slice of all unit conversions
// in the database
func GetConversions(db *sql.DB) ([]Conversion, error) {

	sqlString := "SELECT id, fromUnit, toUnit, multiplicand, denominator, fromOffset, toOffset " +
		"FROM unitConversions WHERE fromUnit = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversions []Conversion

	for rows.Next() {

		var (
			id           int
			fromUnit     int //fk
			toUnit       int //fk
			multiplicand float64
			denominator  float64
			fromOffset   float64
			toOffset     float64
		)
		err = rows.Scan(&id, &fromUnit, &toUnit, &multiplicand, &denominator,
			&fromOffset, &toOffset)
		if err != nil {
			return nil, err
		}
		debugLogger.Printf("ConversionID: %d", id)
		tempConversion := Conversion{ID: id, FromUnit: fromUnit, ToUnit: toUnit,
			Multiplicand: multiplicand, Denominator: denominator,
			FromOffset: fromOffset, ToOffset: toOffset}

		conversions = append(conversions, tempConversion)
	}
	return conversions, nil
}

// GetConversionsToUnit returns a slice of Conversion objects that
// are converting to toUnit
func GetConversionsToUnit(db *sql.DB, toUnit int) ([]Conversion, error) {

	sqlString := "SELECT id, fromUnit, toUnit, multiplicand, denominator, fromOffset, toOffset " +
		"FROM unitConversions WHERE fromUnit = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, toUnit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversions []Conversion

	for rows.Next() {

		var (
			id           int
			fromUnit     int //fk
			toUnit       int //fk
			multiplicand float64
			denominator  float64
			fromOffset   float64
			toOffset     float64
		)
		err = rows.Scan(&id, &fromUnit, &toUnit, &multiplicand, &denominator,
			&fromOffset, &toOffset)
		if err != nil {
			return nil, err
		}
		debugLogger.Printf("ConversionID: %d", id)
		tempConversion := Conversion{ID: id, FromUnit: fromUnit, ToUnit: toUnit,
			Multiplicand: multiplicand, Denominator: denominator,
			FromOffset: fromOffset, ToOffset: toOffset}

		conversions = append(conversions, tempConversion)
	}
	return conversions, nil
}

// GetConversionsFromUnit returns a slice of Conversion objects that
// are converting from: fromUnit
func GetConversionsFromUnit(db *sql.DB, fromUnit int) ([]Conversion, error) {

	sqlString := "SELECT id, fromUnit, toUnit, multiplicand, denominator, fromOffset, toOffset " +
		"FROM unitConversions WHERE fromUnit = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, fromUnit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversions []Conversion

	for rows.Next() {

		var (
			id           int
			fromUnit     int //fk
			toUnit       int //fk
			multiplicand float64
			denominator  float64
			fromOffset   float64
			toOffset     float64
		)
		err = rows.Scan(&id, &fromUnit, &toUnit, &multiplicand, &denominator,
			&fromOffset, &toOffset)
		if err != nil {
			return nil, err
		}
		debugLogger.Printf("ConversionID: %d", id)
		tempConversion := Conversion{ID: id, FromUnit: fromUnit, ToUnit: toUnit,
			Multiplicand: multiplicand, Denominator: denominator,
			FromOffset: fromOffset, ToOffset: toOffset}
		conversions = append(conversions, tempConversion)
	}
	return conversions, nil
}

// GetInventory returns a slice of all Inventory objects in the database
//
//This is used for getting all items in inventory
func GetInventory(db *sql.DB) ([]InventoryItem, error) {

	sqlString := "SELECT id, ean, name, description, quantity, quantityUnit FROM inventory"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventory []InventoryItem

	units, err := GetUnits(db)
	if err != nil {
		return nil, err
	}
	sort.Stable(ByIDU(units)) // sort units slice stably

	for rows.Next() {
		var (
			id          int
			ean         string
			name        string
			description string
			quantity    float64
			unit        int //fk
		)

		err = rows.Scan(&id, &ean, &name, &description, &quantity, &unit)
		if err != nil {
			return nil, err
		}

		// https://golang.org/pkg/sort/#Search
		// first find index
		idIdx := sort.Search(len(units), func(i int) bool { return units[i].ID >= unit })

		// then perform sanity checks
		var tempUnit Unit
		if idIdx < len(units) && units[idIdx].ID == unit {
			tempUnit = units[idIdx]
		} else {
			return nil, fmt.Errorf("ID: %d not found in units slice %v", unit, units)
		}

		debugLogger.Printf("Inventory Item ID: %d, Inventory Item Name: %s", id, name)
		tempInventoryItem := InventoryItem{ID: id, EAN: ean, Name: name,
			Description: description, Quantity: quantity, QuantityUnit: tempUnit}
		inventory = append(inventory, tempInventoryItem)

	}

	return inventory, nil

}

// GetInventoryForIngredient returns a slice of inventory objects associated with
// the specified ingredientID
func GetInventoryForIngredient(db *sql.DB, ingredientID int) ([]InventoryItem, error) {

	sqlString := "SELECT id, ean, name, description, quantity, quantityUnit FROM inventory AS inv " +
		"INNER JOIN ingredient_inventory AS ii ON inv.id = ii.inventoryID WHERE ii.ingredientID = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, ingredientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventory []InventoryItem

	units, err := GetUnits(db)
	if err != nil {
		return nil, err
	}
	sort.Stable(ByIDU(units)) // sort units slice stably

	for rows.Next() {
		var (
			id          int
			ean         string
			name        string
			description string
			quantity    float64
			unit        int //fk
		)

		err = rows.Scan(&id, &ean, &name, &description, &quantity, &unit)
		if err != nil {
			return nil, err
		}
		// https://golang.org/pkg/sort/#Search
		// first find index
		idIdx := sort.Search(len(units), func(i int) bool { return units[i].ID >= unit })

		// then perform sanity checks
		var tempUnit Unit
		if idIdx < len(units) && units[idIdx].ID == unit {
			tempUnit = units[idIdx]
		} else {
			return nil, fmt.Errorf("ID: %d not found in units slice %v", unit, units)
		}

		debugLogger.Printf("Inventory ID: %d, Inventory Item Name: %s", id, name)

		tempInventory := InventoryItem{ID: id, EAN: ean, Name: name, Description: description,
			Quantity: quantity, QuantityUnits: tempUnit}

		inventory = append(inventory, tempInventory)
	}

	return inventory, nil

}

// GetUnits returns a slice of all units  defined in the database
func GetUnits(db *sql.DB) ([]Unit, error) {

	sqlString := "SELECT id, name, description, symbol, isCustom, refIngredient, unitType FROM units"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []Unit

	unitTypes, err := GetUnitTypes(db)
	if err != nil {
		return nil, err
	}

	sort.Stable(ByIDUT(unitTypes))
	ingredients, err := GetIngredients(db)
	if err != nil {
		return nil, err
	}

	sort.Stable(ByIDI(ingredients))

	for rows.Next() {
		var (
			id            int
			name          string
			description   string
			symbol        string
			isCustom      bool
			refIngredient int //fk
			unitType      int //fk
		)

		err = rows.Scan(&id, &name, &description, &symbol, &isCustom, &refIngredient,
			&unitType)
		if err != nil {
			return nil, err
		}

		// https://golang.org/pkg/sort/#Search
		// first find index
		idIdx := sort.Search(len(ingredients), func(i int) bool { return ingredients[i].ID >= refIngredient })

		// then perform sanity checks
		var tempIngredient Ingredient
		if idIdx < len(ingredients) && ingredients[idIdx].ID == refIngredient {
			tempIngredient = ingredients[idIdx]
		} else {
			return nil, fmt.Errorf("ID: %d not found in ingredients slice %v", refIngredient, ingredients)
		}

		// now search unit types
		idIdx = sort.Search(len(unitTypes), func(i int) bool { return unitTypes[i].ID >= unitType })

		// then perform sanity checks
		var tempType UnitType
		if idIdx < len(unitTypes) && unitTypes[idIdx].ID == unitType {
			tempType = unitTypes[idIdx]
		} else {
			return nil, fmt.Errorf("ID: %d not found in ingredients slice %v", unitType, unitTypes)
		}

		debugLogger.Printf("Unit ID: %d, Unit Name: %s", id, name)
		tempUnit := Unit{ID: id, Name: name, Description: description,
			RefIngredient: tempIngredient, UnitType: tempType}
		units = append(units, tempUnit)
	}

	return units, nil
}

// GetUnitTypes returns a slice of UnitType objects that are
// in the database
//
//Note: this table is essentially static
func GetUnitTypes(db *sql.DB) ([]UnitType, error) {

	sqlString := "SELECT id, name FROM unitType"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unitTypes []UnitType

	for rows.Next() {
		var (
			id   int
			name string
		)

		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		debugLogger.Printf("UnitType ID: %d, UnitType Name: %s", id, name)
		tempUnitType := UnitType{ID: id, Name: name}
		unitTypes = append(unitTypes, tempUnitType)
	}
	return unitTypes, nil

}

// GetTagsForRecipe returns a slice of Tag objects containing the tags associated
// with the passed in recipe ID
func GetTagsForRecipe(db *sql.DB, recipeID int) ([]Tag, error) {

	sqlString := "SELECT id, name, description FROM tags INNER JOIN tag_recipe AS tr" +
		"ON tags.id = tr.tagID WHERE tr.recipeID = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag

	for rows.Next() {
		var (
			id          int
			name        string
			description string
		)

		err = rows.Scan(&id, &name, &description)
		if err != nil {
			return nil, err
		}

		debugLogger.Printf("Tag ID: %d, Tag: %s", id, name)
		tempTag := Tag{ID: id, Name: name, Description: description}
		tags = append(tags, tempTag)

	}

	return tags, nil
}

// GetTags returns a slice of Tag objects containing all tags in database
//
// Used to populate a list of tags. Not for getting all attributes of tags
func GetTags(db *sql.DB) ([]Tag, error) {
	sqlString := "SELECT id, name, description FROM tags"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag

	for rows.Next() {
		var (
			id          int
			name        string
			description string
		)

		err = rows.Scan(&id, &name, &description)
		if err != nil {
			return nil, err
		}

		debugLogger.Printf("Tag ID: %d, Tag: %s", id, name)
		tempTag := Tag{ID: id, Name: name, Description: description}
		tags = append(tags, tempTag)
	}

	return tags, nil

}

// IngredientCount returns the number of ingredients in a recipe
//
// Used to size the ingredients table
func IngredientCount(db *sql.DB, recipeID int) (int, error) {

	sqlString := "SELECT count(*) FROM ingredients AS ing INNER JOIN ingredient_recipe AS ir " +
		"ON ing.id = ir.ingredientID WHERE ir.RecipeID = ?"

	ingCount := 0

	debugLogger.Println(sqlString)
	row := db.QueryRow(sqlString, recipeID)
	//defer row.Close()
	err := row.Scan(&ingCount)
	if err != nil {
		return ingCount, err
	}

	return ingCount, nil

}
