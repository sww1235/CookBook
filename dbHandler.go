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

	needInit := false
	requiredTables := []string{"recipes", "ingredients", "ingredient_inventory", "ingredient_recipe",
		"steps", "stepType", "step_recipe", "inventory", "units", "tags", "tag_recipe"}

	for _, table := range requiredTables {
		sqlStatement := fmt.Sprintf("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='%s'", table)
		infoLogger.Println(sqlStatement)
		rows, err := db.Query(sqlStatement)
		if err != nil {
			fatalLogger.Panicln("could not check if table exists", err)
		}
		var rowCount int
		for rows.Next() {

			var count int
			err = rows.Scan(&count)
			if err != nil {
				fatalLogger.Fatalln("reading row count failed", err)
			}
			infoLogger.Println(count)
			rowCount = count
		}

		rows.Close()

		if rowCount == 0 {
			// Assuming that if one table is missing, all tables will need to be recreated.
			needInit = true
			break
		}
	}

	if needInit {
		// need to create tables
		// using map to allow for easier iteration
		createQueries := make(map[string]string)
		createQueries["UnitsTable"] = "CREATE TABLE units(id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, description TEXT)"

		createQueries["RecTable"] = "CREATE TABLE recipes (id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, description TEXT, comments TEXT, source TEXT, author TEXT, " +
			"quantity NUM, quantityUnits INTEGER, FOREIGN KEY(quantityUnits) REFERENCES units(id))"

		createQueries["InvTable"] = "CREATE TABLE inventory (id INTEGER NOT NULL PRIMARY KEY, " +
			"EAN TEXT UNIQUE, name TEXT, description TEXT, quantity NUM, packageQuantity NUM, " +
			"packageQuantityUnits INTEGER, FOREIGN KEY(packageQuantityUnits) REFERENCES units(id))"

		createQueries["IngTable"] = "CREATE TABLE ingredients (id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, quantity NUM, quantityUnits INTEGER, inventoryID INTEGER, " +
			"FOREIGN KEY(inventoryID) REFERENCES inventory(id), " +
			"FOREIGN KEY(quantityUnits) REFERENCES units(id))"

		createQueries["IngInvTable"] = "CREATE TABLE ingredient_inventory( " +
			"ingredientID INTEGER NOT NULL, inventoryID INTEGER NOT NULL, " +
			"FOREIGN KEY(ingredientID) REFERENCES ingredients(id), " +
			"FOREIGN KEY(inventoryID) REFERENCES inventory(id), " +
			"PRIMARY KEY(ingredientID, inventoryID))"

		createQueries["IngRecTable"] = "CREATE TABLE ingredient_recipe( " +
			"ingredientID INTEGER NOT NULL, recipeID INTEGER NOT NULL, " +
			"FOREIGN KEY(ingredientID) REFERENCES ingredients(id), " +
			"FOREIGN KEY(recipeID) REFERENCES recipes(id), " +
			"PRIMARY KEY(ingredientID, recipeID))"

		createQueries["StepTypeTable"] = "CREATE TABLE stepType (id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT)"

		createQueries["StepTable"] = "CREATE TABLE steps( id INTEGER NOT NULL PRIMARY KEY, " +
			"instructions TEXT, time NUM, stepTypeID INTEGER, temperature NUM, tempUnits INTEGER, " +
			"FOREIGN KEY(stepTypeID) REFERENCES stepType(id), " +
			"FOREIGN KEY(tempUnits) REFERENCES units(id))"

		createQueries["StepRecTable"] = "CREATE TABLE step_recipe( stepID INTEGER NOT NULL, " +
			"recipeID INTEGER NOT NULL, " +
			"FOREIGN KEY(stepID) REFERENCES steps(id), " +
			"FOREIGN KEY(recipeID) REFERENCES recipes(id), " +
			"PRIMARY KEY(stepID, recipeID))"
		createQueries["TagTable"] = "CREATE TABLE tags(id INTEGER NOT NULL PRIMARY KEY, name TEXT NOT NULL)"

		createQueries["TagRecTable"] = "CREATE TABLE tag_recipe( tagID INTEGER NOT NULL, recipeID INTEGER NOT NULL, " +
			"FOREIGN KEY(tagID) REFERENCES tags(id), " +
			"FOREIGN KEY(recipeID) REFERENCES recipes(id), " +
			"PRIMARY KEY(tagID, recipeID))"

		// since not all tables exist, for now drop all tables, then recreate them

		// now create all the tables
		for table, query := range createQueries {

			_, err := db.Exec(query)
			if err != nil {
				fatalLogger.Fatalf("Failed to create table: %s due to error: %s", table, err)
			}
		}

	}
	return db

}

func insertRecipe(db *sql.DB, recipe backend.Recipe) error {

	return nil
}
