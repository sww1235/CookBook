package main

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

//initDB opens a connection to a sqlite database that stores recipes
func initDB(databasePath string) *sql.DB {

	// TODO: allow for other database backends, and allow for tables to be in separate backends

	// try to create database directory, then check if database exists before opening.
	// race condition exists on check if file exists, but I don't think this application
	// will run into it.

	mkErr := os.MkdirAll(path.Dir(databasePath), 0744)
	if mkErr != nil {
		fatalLogger.Panicf("Unable to use %s as recipe directory, Err: %s",
			path.Dir(databasePath), mkErr)
	}

	_, err := os.Stat(databasePath)
	if os.IsNotExist(err) {
		infoLogger.Printf("database doesn't exist, creating now at path %s", databasePath)
	}

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		fatalLogger.Panicln("Could not open recipe database", err)
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
		"stepType":             false,
		"step_recipe":          false,
		"inventory":            false,
		"units":                false,
		"tags":                 false,
		"tag_recipe":           false,
		"lastMade":             false,
		"equipment":            false,
		"equipment_recipe":     false,
		"unitType":             false,
		"unitConversions":      false,
	}

	for table := range requiredTables {
		sqlStatement := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"

		var rowCount int
		err := db.QueryRow(sqlStatement, table).Scan(&rowCount)
		if err != nil {
			fatalLogger.Panicln("could not check if table exists", err)
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
		fatalLogger.Panicln("Existing database missing at least one critical table. See log messages above." +
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

		createQueries["stepType"] = "CREATE TABLE stepType (id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT)"

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
				fatalLogger.Panicf("Failed to create table: %s due to error: %s", table, err)
			}
		}

	}
	return db

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

// GetRecipes returns a map of recipe names indexed on their database id
//
// Used to populate a list of recipes. Not for getting all attributes of recipes
func GetRecipes(db *sql.DB) (map[int]string, error) {

	sqlString := "SELECT id, name FROM recipes"

	listRecipes := make(map[int]string)

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return listRecipes, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			fatalLogger.Panicln("reading row failed", err)
		}
		debugLogger.Printf("Recipe ID: %d, Recipe Name: %s", id, name)
		// add row values to recipe map
		listRecipes[id] = name
	}

	return listRecipes, nil

}

// GetTags returns a map of tag names indexed on their database id
//
// Used to populate a list of tags. Not for getting all attributes of tags
func GetTags(db *sql.DB) (map[int]string, error) {

	sqlString := "SELECT id, name FROM tags"

	listTags := make(map[int]string)

	debugLogger.Println(sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return listTags, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			fatalLogger.Panicln("reading row failed", err)
		}
		debugLogger.Printf("Tag ID: %d, Tag Name: %s", id, name)
		// add row values to recipe map
		listTags[id] = name
	}

	return listTags, nil

}
