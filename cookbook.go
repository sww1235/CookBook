package main

//TODO: either split units out of ingredient.go and into their own file, or at
//least implement a way to add conversions between volume units (ie: cups to
//teespoons) and weight units (grams to lbs), as well as be able to add custom
//conversions between a specific weight unit and a specific volume unit for a
//specific ingredient

// TODO: implement config settings for forground and background colors for cui

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/awesome-gocui/gocui"
	backend "github.com/sww1235/recipe-database"
)

//Configuration stores the configuration that is read in and out from a file

var viewedRecipe string
var addRecipeToggle bool
var httpServer bool
var httpServerFlagIP string

var config Configuration

var infoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
var fatalLogger = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	// read in config file and command line flags
	err := initialization()
	if err != nil {
		fatalLogger.Panicln("Config and flag init failed", err)
	}

	db := initDB(config.RecipeDatabase)

	if viewedRecipe != "" {
		err := displaySingleRecipe(viewedRecipe)
		if err != nil {
			fatalLogger.Fatalf("%s not found in Recipes, check your spelling and capitilization\n",
				viewedRecipe)
		}
		os.Exit(0)
	} else if addRecipeToggle {
		//read in recipe from commandline
		tempRecipe, err := backend.ReadRecipe()
		if err != nil {
			fatalLogger.Fatalln("Error reading recipe from command line:", err)
		}
		//insert it into database

		err = insertRecipe(db, tempRecipe)
		if err != nil {
			fatalLogger.Fatalln("Error inserting new recipe into database:", err)
		}
		os.Exit(0)
	} else if httpServer {
		err := startHTTPServer()
		if err != nil {
			fatalLogger.Fatalln("Something went wrong with the HTTP server", err)
		}
	} else {
		err := startCUI()
		if err != nil && err != gocui.ErrQuit {
			fatalLogger.Fatalln("Something went wrong with the CUI", err)
		}
	}

}

//initialization sets up command line options, logging and config stuffs
//read config file, either from -c flag or default in ~/.config
func initialization() error {

	//if this fails, cannot create default recipeDir or configDir. This is fatal
	currUserHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	currUserConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	//default paths. Will not be overridden
	defaultConfigPath := path.Join(currUserConfigDir, "cookbook.cfg")

	defaultRecipeDatabaseDir := path.Join(currUserHomeDir, "cookbook")

	//save default paths
	defaultRecipeDatabase := path.Join(defaultRecipeDatabaseDir, "cookbook.db")
	const defaultServerIP = "127.0.0.1"

	//setting config to defaults
	config.RecipeDatabase = defaultRecipeDatabase

	//Define and Parse commandline flags here
	//Defaults are set in flags as appropriate
	flagConfigPath := flag.String("c", defaultConfigPath, "Path to config file")
	const flagViewedRecipeUsage = "Recipe to view. Recipe name is case sensitive " +
		"and must be typed exactly. This flag is provided as a courtesy for " +
		"scripting and people who can't run termbox"
	flagViewedRecipe := flag.String("v", "", flagViewedRecipeUsage)
	flagRecipeDatabaseDir := flag.String("r", defaultRecipeDatabaseDir, "Directory to store recipe database in")
	flagAddRecipeToggle := flag.Bool("n", false, "Add new recipe")
	flagHTTPServer := flag.Bool("H", false, "Use HTTP server instead of terminal")
	flagIPConfig := flag.String("ip", defaultServerIP, "IP to start HTTP server on")
	flag.Parse()

	//Retrieve some values from flags and set global variables
	//if these flags are not set, defaults will be set
	viewedRecipe = *flagViewedRecipe
	addRecipeToggle = *flagAddRecipeToggle
	httpServer = *flagHTTPServer
	httpServerFlagIP = *flagIPConfig

	if *flagConfigPath != defaultConfigPath {
		infoLogger.Println("Using config file path from flag", *flagConfigPath)

	}

	//Attempt to read config
	config, cfgErr := readConfig(*flagConfigPath)
	if cfgErr != nil {
		infoLogger.Printf("config file not openable at path: %s. Err: %s. "+
			"Using default configuration", *flagConfigPath, cfgErr)
	}

	//try flag recipeDatabaseDir
	if *flagRecipeDatabaseDir != defaultRecipeDatabaseDir {
		infoLogger.Printf("Using recipeDatabaseDir %s from flag", *flagRecipeDatabaseDir)
		config.RecipeDatabase = path.Join(*flagRecipeDatabaseDir, "cookbook.db")
	} else if path.Dir(config.RecipeDatabase) != defaultRecipeDatabaseDir {
		infoLogger.Printf("Using recipeDatabaseDir %s from config", path.Dir(config.RecipeDatabase))
		//config.RecipeDatabase already set
	} else {
		infoLogger.Printf("Using default recipeDir %s", defaultRecipeDatabaseDir)
		config.RecipeDatabase = defaultRecipeDatabase
	}
	return nil
}

func finalize() {

	// close database
	//

}

//check for and read ingredient database

//establish ncurses-like gui, look at termbox-go or gocui, mop-tracker
//or hecate for example programs

//add new recipe function

//displaySingleRecipe prints a full recipe to stdout using the recipe.print method
//recipeName is passed into sql prepared statement.
//Multiple recipes can be returned from sql query, and so the user is prompted
//for which one they want.
func displaySingleRecipe(recipeName string) error {
	var tempRecipe backend.Recipe
	fmt.Print(tempRecipe.String())

	fmt.Println("Press enter to exit program")
	//will keep attempting to read from stdin until it receives a '\n'
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

//view recipe function

//format recipe to markdown

//Functions to shutdown program
//-write out all new and changed recipes to files
