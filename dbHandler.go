package main

import (
	"database/sql"
	"fmt"
	"os"
	"path"

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

		createQueries["steps"] = "CREATE TABLE steps( id INTEGER NOT NULL PRIMARY KEY, " +
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
		_, err := InsertIngredient(db, ing, r, ing.InventoryReference.ID)

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

	insertSQL := "INSERT INTO steps (instructions, time, stepType, temperature, tempUnit) " +
		"VALUES(?, ?, ?, ?, ?)"

	mapSQL := "INSERT INTO step_recipe (recipeID, stepID) VALUES (?, ?)"

	// if ingredient already exists in database
	if stp.ID != -1 {
		return -1, fmt.Errorf("Step %s already exists in database", stp)
	}

	result, err := db.Exec(insertSQL, stp.Instructions, stp.TimeNeeded, stp.StepType,
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

	updateSQL := "UPDATE steps SET instructions = ?, time = ?, stepType = ? " +
		"temperature = ?, tempUnit = ? WHERE id = ?"

	// if ingredient doesn't exist in database
	if stp.ID == -1 {
		return -1, fmt.Errorf("Step %s doesn't exist in database, can't update it", stp)

	}

	_, err := db.Exec(updateSQL, stp.Instructions, stp.TimeNeeded, stp.StepType,
		stp.Temperature.Value, stp.Temperature.Unit.ID, stp.ID)

	//TODO: need to delete and remap ingredient_inventory

	return stp.ID, err
}

// GetRecipes returns a map of Recipe structs indexed on
// the database table index
func GetRecipes(db *sql.DB) (map[int]Recipe, error) {

	sqlString := "SELECT * FROM recipes"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipes := make(map[int]Recipe)
	units, err := GetUnit(db)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			id             int
			name           string
			description    string
			comments       string
			source         string
			author         string
			quantity       float64
			qtyUnitsid     int //fk
			initialVersion int //fk
			version        int
		)

		err = rows.Scan(&id, &name, &description, &comments, &source, &author,
			&quantity, &qtyUnitsid, &initialVersion, &version)
		if err != nil {
			return nil, err
		}
		ingredients, err := GetIngredinetsForRecipe(db, id)
		if err != nil {
			return nil, err
		}
		steps, err := GetStepsForRecipe(db, id)
		if err != nil {
			return nil, err
		}
		equipment, err := GetEquipmentForRecipe(db, id)
		if err != nil {
			return nil, err
		}
		tags, err := GetTagsForRecipe(db, id)
		if err != nil {
			return nil, err
		}
		debugLogger.Printf("Recipe ID: %d, Recipe Name: %s", id, name)
		// add row values to recipe map
		recipes[id] = Recipe{Id: id, Name: name, Description: description, Comments: comments,
			Source: source, Author: author, Ingredients: ingredients, QuantityMade: quantity,
			QuantityMadeUnits: units[qtyUnitsid], Steps: steps,
			Equipment: equipment, Tags: tags, Version: version, InitialVersion: initialVersion}
	}

	return recipes, nil

}

// GetEquipmentForRecipe returns a map of Equipment structs indexed on
// the database table index for the specified recipe
func GetEquipmentForRecipe(db *sql.DB, recipeID int) (map[int]Equipment, error) {

	sqlString := "SELECT * FROM equipment AS equip INNER JOIN equipment_recipe AS er" +
		"ON equip.id = er.equipmentID WHERE er.recipeID = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	equipment := make(map[int]Equipment)

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
		equipment[id] = Equipment{Id: id, Name: name, IsOwned: isOwned}
	}

	return equipment, nil
}

// GetIngredientsForRecipe returns a map of Ingredient structs indexed on
// the database table index for the specified recipe
func GetIngredientsForRecipe(db *sql.DB, recipeID int) (map[int]Ingredient, error) {

	sqlString := "SELECT * FROM ingredients AS ing INNER JOIN ingredient_recipe AS ir" +
		"ON ing.id = ir.ingredientID WHERE ir.recipeID = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ingredients := make(map[int]Ingredient)

	units, err := GetUnit(db)
	if err != nil {
		return nil, err
	}

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
		inventory, err := GetInventoryForIngredient(db, id)
		if err != nil {
			return nil, err
		}

		var invStruct []Inventory
		for _, inv := range inventory {
			invStruct = append(invStruct, inv)
		}

		debugLogger.Printf("Ingredient ID: %d, Ingredient Name: %s", id, name)
		ingredients[id] = Ingredient{Id: id, Name: name, QuantityUsed: quantity,
			QuantityUsedUnits: units[unit], InventoryReference: invStruct}
	}

	return ingredients, nil
}

// GetStepsForRecipe returns a map of Step structs indexed on
// the database table index for the specified recipe
func GetStepsForRecipe(db *sql.DB, recipeID int) (map[int]Step, error) {

	sqlString := "SELECT * FROM steps INNER JOIN step_recipe AS is" +
		"ON steps.id = is.stepID WHERE ir.recipeID = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	steps := make(map[int]Step)

	units, err := GetUnit(db)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id           int
			instructions string
			timeNeeded   int //seconds
			stepType     int //fk to stepType iota
			temperature  float64
			tempUnit     int //fk
		)

		err = rows.Scan(&id, &instructions, &timeNeeded, &stepType,
			&temperature, &tempUnit)
		if err != nil {
			return nil, err
		}

		duration := time.ParseDuration(strconv.Itoa(timeNeeded) + s)
		temperatureUnit := units[tempUnit]
		temp := Temperature{Value: temperature, Unit: temperatureUnit}

		debugLogger.Printf("Step ID: %d", id)
		steps[id] = Step{Id: id, TimeNeeded: duration, StepType: StepType(stepType),
			Temperature: temp}
	}

	return ingredients, nil
}

// GetIngredients returns a map of Ingredient structs indexed on
// the database table index

func GetIngredients(db *sql.DB) (map[int]Ingredient, error) {

	sqlString := "SELECT * FROM ingredients"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ingredients := make(map[int]Ingredient)

	units, err := GetUnit(db)
	if err != nil {
		return nil, err
	}

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
		inventory, err := GetInventoryForIngredient(db, id)
		if err != nil {
			return nil, err
		}

		var invStruct []Inventory
		for _, inv := range inventory {
			invStruct = append(invStruct, inv)
		}

		debugLogger.Printf("Ingredient ID: %d, Ingredient Name: %s", id, name)
		ingredients[id] = Ingredient{Id: id, Name: name, QuantityUsed: quantity,
			QuantityUsedUnits: units[unit], InventoryReference: invStruct}
	}

	return ingredients, nil
}

// GetConversions returns a map of Conversion objects indexed on
// the database table index
func GetConversions(db *sql.DB) (map[int]Conversion, error) {

	sqlString := "SELECT * FROM unitConversions"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversions := make(map[int]Conversion, error)

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
		conversions[id] = Conversion{Id: id, FromUnit: fromUnit, ToUnit: toUnit,
			Multiplicand: multiplicand, Denominator: denominator,
			FromOffset: fromOffset, ToOffset: toOffset}

	}
}

// GetConversionsToUnit returns a map of Conversion objects indexed on
// the database table index toUnit
func GetConversionsToUnit(db *sql.DB, toUnit int) (map[int]Conversion, error) {

	sqlString := "SELECT * FROM unitConversions WHERE toUnit = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, toUnit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversions := make(map[int]Conversion, error)

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
		conversions[id] = Conversion{Id: id, FromUnit: fromUnit, ToUnit: toUnit,
			Multiplicand: multiplicand, Denominator: denominator,
			FromOffset: fromOffset, ToOffset: toOffset}

	}
}

// GetConversionsFromUnit returns a map of Conversion objects indexed on
// the database table index fromUnit
func GetConversionsFromUnit(db *sql.DB, fromUnit int) (map[int]Conversion, error) {

	sqlString := "SELECT * FROM unitConversions WHERE fromUnit = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, fromUnit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversions := make(map[int]Conversion, error)

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
		conversions[id] = Conversion{Id: id, FromUnit: fromUnit, ToUnit: toUnit,
			Multiplicand: multiplicand, Denominator: denominator,
			FromOffset: fromOffset, ToOffset: toOffset}

	}
}

// GetInventory returns a map of Inventory structs indexed on
// the database table index
//
//This is used for getting all items in inventory
func GetInventory(db *sql.DB) (map[int]Inventory, error) {

	sqlString := "SELECT * FROM inventory"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventory := make(map[int]Inventory)

	units, err := GetUnit(db)
	if err != nil {
		return nil, err
	}

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
		debugLogger.Printf("Inventory ID: %d, Inventory Item Name: %s", id, name)
		inventory[id] = Inventory{Id: id, EAN: ean, Name: name, Description: description, Quantity: quantity,
			QuantityUnits: units[unit]}
	}

	return ingredients, nil

}

// GetInventoryForIngredient returns a map of Inventory structs
// indexed on the database table index associated with
// the specified ingredientID
func GetInventoryForIngredient(db *sql.DB, ingredientID int) (map[int]Inventory, error) {

	sqlString := "SELECT * FROM inventory AS inv INNER JOIN ingredient_inventory AS ii " +
		"ON inv.id = ii.inventoryID WHERE ii.ingredientID = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, ingredientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventory := make(map[int]Inventory)

	units, err := GetUnit(db)
	if err != nil {
		return nil, err
	}

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
		debugLogger.Printf("Inventory ID: %d, Inventory Item Name: %s", id, name)
		inventory[id] = Inventory{Id: id, EAN: ean, Name: name, Description: description, Quantity: quantity,
			QuantityUnits: units[unit]}
	}

	return ingredients, nil

}

// GetUnits returns a map of Unit structs indexed on
// the database table index
func GetUnits(db *sql.DB) (map[int]Unit, error) {

	sqlString := "SELECT * FROM units"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	units := make(map[int]Unit)

	unitTypes, err := GetUnitTypes(db)
	if err != nil {
		return nil, err
	}

	ingredients, err := GetIngredients(db)
	if err != nil {
		return nil, err
	}

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
		debugLogger.Printf("Unit ID: %d, Unit Name: %s", id, name)
		units[id] = Unit{Id: id, Name: name, Description: description,
			RefIngredient: ingredients[refIngredient], UnitType: unitTypes[unitType]}
	}

	return units, nil
}

// GetUnitTypes returns a map of UnitType structs indexed on
// the database table index.
func GetUnitTypes(db *sql.DB) (map[int]UnitType, error) {

	sqlString := "SELECT * FROM unitType"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	unitTypes := make(map[int]UnitType)

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

		unitTypes[id] = UnitType{Id: id, Name: name}
	}
	return unitTypes, nil

}

// GetTagsForRecipe
func GetTagsForRecipe(db *sql.DB, recipeID int) (map[int]Tag, error) {

	sqlString := "SELECT * FROM tags INNER JOIN tag_recipe AS tr" +
		"ON tags.id = tr.tagID WHERE tr.recipeID = ?"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make(map[int]Tag)

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
		tags[id] = Tag{Id: id, Name: name, Description: description}
	}

	return tags, nil
}

// GetTags returns a map of tag names indexed on their database id
//
// Used to populate a list of tags. Not for getting all attributes of tags
func GetTags(db *sql.DB) (map[int]string, error) {
	sqlString := "SELECT * FROM tags"

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make(map[int]Tag)

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
		tags[id] = Tag{Id: id, Name: name, Description: description}
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
