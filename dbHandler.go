package main

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
	backend "github.com/sww1235/recipe-database"
)

//initDB opens a connection to a sqlite database that stores recipes
func initDB(databasePath string) *sql.DB {

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
	requiredTables := []string{"recipe", "ingredient", "ingredient_inventory", "ingredient_recipe",
		"step", "stepType", "step_recipe", "inventory", "units", "tags", "tag_recipe"}

	for _, table := range requiredTables {
		sqlStatement := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='" + table + "'"
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
		createUnitsTableQuery := "CREATE TABLE units(id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, description TEXT)"
		createRecTableQuery := "CREATE TABLE recipes (id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, description TEXT, comments TEXT, source TEXT, author TEXT, " +
			"quantity NUM, FOREIGN KEY(quantityUnits) REFERENCES units(id))"
		createInvTableQuery := "CREATE TABLE inventory (id INTEGER NOT NULL PRIMARY KEY, " +
			"EAN TEXT UNIQUE, name TEXT, description TEXT, quantity NUM, packageQuantity NUM, " +
			"FOREIGN KEY(packageQuantityUnits) REFERENCES units(id))"
		createIngTableQuery := "CREATE TABLE ingredients (id INTEGER NOT NULL PRIMARY KEY, " +
			"name TEXT, FOREIGN KEY(inventoryID) REFERENCES inventory(id), quantity NUM, " +
			"FOREIGN KEY(quantityUnits) REFERENCES units(id))"
		createIngInvTableQuery := ""
		createStepTableQuery := "CREATE TABLE steps (id INTEGER NOT NULL PRIMARY KEY"
		createStepTypeTableQuery := "CREATE TABLE stepType (id INTEGER NOT NULL PRIMARY KEY"
		createStepRecTableQuery := ""
		createTagTableQuery := ""
		createTagRecTableQuery := ""

	}
	return db

}

func insertRecipe(db *sql.DB, recipe backend.Recipe) error {

	return nil
}
