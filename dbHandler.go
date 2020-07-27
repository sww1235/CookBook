package main

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
	backend "github.com/sww1235/recipe-database"
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
		"unitType":             false,
		"unitConversions":      false,
	}

	for table := range requiredTables {
		sqlStatement := fmt.Sprintf("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='%s'", table)
		debugLogger.Println(sqlStatement)
		rows, err := db.Query(sqlStatement)
		if err != nil {
			fatalLogger.Panicln("could not check if table exists", err)
		}
		var rowCount int
		for rows.Next() {

			var count int
			err = rows.Scan(&count)
			if err != nil {
				fatalLogger.Panicln("reading row count failed", err)
			}
			debugLogger.Println(count)
			rowCount = count
		}

		rows.Close()

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
			infoLogger.Printf("Table %s missing. Manually create this table, or delete database so it can be recreated.", table)
		}
	}

	if missingTable {
		fatalLogger.Panicln("Existing database missing critical table. See log messages above.")
	}

	if needInit {
		// need to create tables
		// using map to allow for easier iteration
		createQueries := make(map[string]string)
		createQueries["units"] = "CREATE TABLE units(id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, description TEXT)"

		createQueries["recipes"] = "CREATE TABLE recipes (id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, description TEXT, comments TEXT, source TEXT, author TEXT, " +
			"quantity NUM, quantityUnits INTEGER, initialVersion INTEGER, version INTEGER, " +
			"FOREIGN KEY(quantityUnits) REFERENCES units(id), " +
			"FOREIGN KEY(initialVersion) REFERENCES recipes(id))"

		createQueries["inventory"] = "CREATE TABLE inventory (id INTEGER NOT NULL PRIMARY KEY, " +
			"EAN TEXT UNIQUE, name TEXT, description TEXT, quantity NUM, packageQuantity NUM, " +
			"packageQuantityUnits INTEGER, FOREIGN KEY(packageQuantityUnits) REFERENCES units(id))"

		createQueries["ingredients"] = "CREATE TABLE ingredients (id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, quantity NUM, quantityUnits INTEGER, inventoryID INTEGER, " +
			"FOREIGN KEY(inventoryID) REFERENCES inventory(id), " +
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

		createQueries["stepType"] = "CREATE TABLE stepType (id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT)"

		createQueries["steps"] = "CREATE TABLE steps( id INTEGER NOT NULL PRIMARY KEY, " +
			"instructions TEXT, time NUM, stepTypeID INTEGER, temperature NUM, tempUnits INTEGER, " +
			"FOREIGN KEY(stepTypeID) REFERENCES stepType(id), " +
			"FOREIGN KEY(tempUnits) REFERENCES units(id))"

		createQueries["step_recipe"] = "CREATE TABLE step_recipe( stepID INTEGER NOT NULL, " +
			"recipeID INTEGER NOT NULL, " +
			"FOREIGN KEY(stepID) REFERENCES steps(id), " +
			"FOREIGN KEY(recipeID) REFERENCES recipes(id), " +
			"PRIMARY KEY(stepID, recipeID))"

		createQueries["tags"] = "CREATE TABLE tags(id INTEGER NOT NULL PRIMARY KEY, name TEXT NOT NULL)"

		createQueries["tag_recipe"] = "CREATE TABLE tag_recipe( tagID INTEGER NOT NULL, recipeID INTEGER NOT NULL, " +
			"FOREIGN KEY(tagID) REFERENCES tags(id), " +
			"FOREIGN KEY(recipeID) REFERENCES recipes(id), " +
			"PRIMARY KEY(tagID, recipeID))"

		createQueries["lastMade"] = "CREATE TABLE lastMade(id INTEGER NOT NULL PRIMARY KEY, " +
			"recipe INTEGER NOT NULL, dateMade TEXT, notes TEXT, " +
			"FOREIGN KEY(recipe) REFERENCES recipes(id))"

		createQueries["equipment"] = "CREATE TABLE equipment(id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, isOwned NUM)"

		createQueries["unitType"] = "CREATE TABLE unitType(id INTEGER NOT NULL PRIMARY KEY, name TEXT)"

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

func insertRecipe(db *sql.DB, recipe backend.Recipe) error {

	return nil
}
