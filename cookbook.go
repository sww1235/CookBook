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
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

//Configuration stores the configuration that is read in and out from a file

var viewedRecipe string
var addRecipeToggle bool
var httpServer bool
var httpServerFlagIP string

var config Configuration

var units map[string]Unit

var debugLogger = log.New(ioutil.Discard, "DEBUG: ", 0)
var infoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
var fatalLogger = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	// read in config file and command line flags
	err := initialization()
	if err != nil {
		fatalLogger.Panicln("Config and flag init failed", err)
	}

	db := initDB(config.RecipeDatabase)

	defer db.Close()

	if viewedRecipe != "" {
		err := displaySingleRecipe(viewedRecipe)
		if err != nil {
			fatalLogger.Panicf("%s not found in Recipes, check your spelling and capitilization\n",
				viewedRecipe)
		}
	} else if addRecipeToggle {
		//read in recipe from commandline
		tempRecipe, err := ReadRecipe()
		if err != nil {
			fatalLogger.Panicln("Error reading recipe from command line:", err)
		}
		//insert it into database

		_, err = InsertRecipe(db, tempRecipe) // don't need recipe id in this case
		if err != nil {
			fatalLogger.Panicln("Error inserting new recipe into database:", err)
		}
	} else if httpServer {
		err := startHTTPServer()
		if err != nil {
			fatalLogger.Panicln("Something went wrong with the HTTP server", err)
		}
	} else {
		err := startCUI()
		if err != nil {
			fatalLogger.Panicln("Something went wrong with the CUI", err)
		}
	}

	// defers in main will run here
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
	flagViewedRecipe := flag.String("r", "", flagViewedRecipeUsage)
	flagRecipeDatabaseDir := flag.String("db", defaultRecipeDatabaseDir, "Directory to store recipe database")
	flagAddRecipeToggle := flag.Bool("n", false, "Add new recipe")
	flagHTTPServer := flag.Bool("H", false, "Use HTTP server instead of terminal")
	flagIPConfig := flag.String("ip", defaultServerIP, "IP to start HTTP server on")
	flagDebugLogging := flag.Bool("D", false, "Show debug logs")
	flag.Parse()

	if *flagDebugLogging {
		debugLogger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		debugLogger.SetOutput(os.Stdout)
	}

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
	tempConfig, cfgErr := readConfig(*flagConfigPath)
	if cfgErr != nil {
		infoLogger.Printf("config file not openable at path: %s. Err: %s. "+
			"Using default configuration", *flagConfigPath, cfgErr)
	} else {
		config = tempConfig
	}

	//try flag recipeDatabaseDir
	if *flagRecipeDatabaseDir != defaultRecipeDatabaseDir { // flag is present if value isn't default
		infoLogger.Printf("Using recipeDatabaseDir %s from flag", *flagRecipeDatabaseDir)
		config.RecipeDatabase = path.Join(*flagRecipeDatabaseDir, "cookbook.db")

		// if config set path to different path than default
	} else if config.RecipeDatabase != defaultRecipeDatabase {
		infoLogger.Printf("!! config value: %s, default value: %s", config.RecipeDatabase, defaultRecipeDatabase)
		infoLogger.Printf("Using recipeDatabaseDir %s from config", path.Dir(config.RecipeDatabase))
		//config.RecipeDatabase already set
	} else {
		// config.RecipeDatabase already set at line 105
		infoLogger.Printf("Using default recipe database %s", defaultRecipeDatabase)
	}
	return nil
}

//ReadRecipe creates a recipe struct by prompting the user for input
func ReadRecipe() (Recipe, error) {

	var tempRecipe Recipe
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Adding new recipe to database. Press Ctrl-C to abort")
	fmt.Print("Enter recipe Name: ")
	tempString, err := reader.ReadString('\n')
	if err != nil {
		return tempRecipe, err
	}
	tempRecipe.Name = tempString

	// read in quantity made and units
	fmt.Print("Enter the number of things the recipe makes: ")

	tempString, err = reader.ReadString('\n')
	if err != nil {
		return tempRecipe, err
	}
	tempRecipe.QuantityMade, err = strconv.Atoi(tempString)
	if err != nil {
		return tempRecipe, err
	}

	for {
		tempString, err = reader.ReadString('\n')
		if err != nil {
			return tempRecipe, err
		}

		tempUnit, ok := units[tempString]

		if ok {

			tempRecipe.QuantityMadeUnit = tempUnit
			break
		}
	}

	// read in ingredients

	// read in steps

	return tempRecipe, nil
}

//displaySingleRecipe prints a full recipe to stdout using the recipe.print method
//recipeName is passed into sql prepared statement.
//Multiple recipes can be returned from sql query, and so the user is prompted
//for which one they want.
func displaySingleRecipe(recipeName string) error {
	var tempRecipe Recipe
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
