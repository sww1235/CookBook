package main

//TODO: either split units out of ingredient.go and into their own file, or at
//least implement a way to add conversions between volume units (ie: cups to
//teespoons) and weight units (grams to lbs), as well as be able to add custom
//conversions between a specific weight unit and a specific volume unit for a
//specific ingredient

//TODO: figure out logging

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"

	backend "github.com/sww1235/recipe-database-backend"
)

//Configuration stores the configuration that is read in and out from a file

var viewedRecipe string
var addRecipeToggle bool
var httpServer bool
var httpServerFlagIP string

var config Configuration

var recipes []backend.Recipe

var infoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
var fatalLogger = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	err := initialization()
	if err != nil {
		fatalLogger.Panicln("Serious issue detected", err)
	}

	readRecipes(config.RecipeDir)

	if viewedRecipe != "" {
		err := displaySingleRecipe(viewedRecipe)
		if err != nil {
			fatalLogger.Fatalf("%s not found in Recipes, check your spelling and capitilization", viewedRecipe)
		}
		os.Exit(0)
	}
	if addRecipeToggle {
		//add recipe needs to be fixed
		recipes = append(recipes, addRecipe())
		saveRecipes(recipes)
		// if config file found then write config
		if !config.configFileNotFound {
			writeConfig(config, config.configPath)
		}

		os.Exit(0)
	}
	if httpServer {
		err := startHTTPServer()
		if err != nil {
			fatalLogger.Fatalln("Something went wrong with the HTTP server", err)
		}
		saveRecipes(recipes)
		// if config file found then write config
		if !config.configFileNotFound {
			writeConfig(config, config.configPath)
		}
	}

	saveRecipes(recipes)
	// if config file found then write config
	if !config.configFileNotFound {
		writeConfig(config, config.configPath)
	}

}

//initialization sets up command line options, logging and config stuffs
//read config file, either from -c flag or default in ~/.config
func initialization() error {

	//if this fails, cannot create default recipeDir or configDir. This is fatal
	currUsr, usrErr := user.Current()
	if usrErr != nil {
		return usrErr
	}
	useDefaltsIfError := func(err bool) {
		config.IPConfig = "127.0.0.1"
		config.RecipeDir = path.Join(currUsr.HomeDir, "cookbook")
		config.configFileNotFound = false
		if err {
			config.configPath = ""
			config.configFileNotFound = true
		}

	}

	//default paths. Will not be overridden
	configDir := path.Join(currUsr.HomeDir, ".config", "cookbook")
	recipeDir := path.Join(currUsr.HomeDir, "cookbook")

	//setting more defaults
	defaultConfigPath := path.Join(configDir, "cookbook.cfg")
	defaultRecipeDir := recipeDir

	//setting config to defaults
	config.configPath = defaultConfigPath
	config.RecipeDir = defaultRecipeDir

	//Define and Parse commandline flags here
	flagConfigPath := flag.String("c", defaultConfigPath, "Path to config file")
	flagViewedRecipe := flag.String("v", "", "Recipe to view. Recipe name is case sensitive and must be typed exactly. This flag is provided as a courtesy for scripting and people who can't run termbox")
	flagRecipeDir := flag.String("r", defaultRecipeDir, "Directory to store recipes in")
	flagAddRecipeToggle := flag.Bool("n", false, "Add new recipe")
	flagHTTPServer := flag.Bool("H", false, "Use HTTP server instead of terminal")
	flagIPConfig := flag.String("ip", "127.0.0.1", "IP to start HTTP server on")
	flag.Parse()

	//Retrieve some values from flags and set global variables
	//if these flags are not set, defaults will be set
	viewedRecipe = *flagViewedRecipe
	addRecipeToggle = *flagAddRecipeToggle
	httpServer = *flagHTTPServer
	httpServerFlagIP = *flagIPConfig

	//Attempt to read config from default location
	defaultConfig, defaultCfgErr := readConfig(config.configPath)

	//when this block is triggered, the config file is not in its default location
	//Or the configPath flag is set
	if *flagConfigPath != defaultConfigPath || defaultCfgErr != nil {
		infoLogger.Println("Config file not found in default location or commandline flag is set, trying commandline flag")

		//config flag is not set and config file does not exist in default location

		if *flagConfigPath == defaultConfigPath {
			//log that config dir is being kept as default due to no flag
			infoLogger.Println("config flag not set, creating directories and config file in default location")
		} else {
			//set configdir to path entered by flag
			infoLogger.Printf("config flag set, using %s as configPath", *flagConfigPath)
			config.configPath = *flagConfigPath
		}
		//Attempt to create configdir in location defined above
		mkErr := os.MkdirAll(filepath.Dir(config.configPath), 0644)

		//MkdirAll returns nil on already exists or dir created, therefore errors are
		//serious, IE directory could not be created due to permissions
		if mkErr != nil {
			infoLogger.Println("Something went wrong attempting to create config directory. Check permissions to specified config directory", mkErr)
			infoLogger.Println("will not attempt to save config file as specified dir cannot be created")
			useDefaltsIfError(true)
		} else {
			readConfig, readCfgErr := readConfig(config.configPath)
			if readCfgErr != nil {
				infoLogger.Println("Config file does not exist at either default or commandline flag location, using defaults", readCfgErr)
				useDefaltsIfError(false)
			} else {
				config.IPConfig = readConfig.IPConfig
				config.RecipeDir = readConfig.RecipeDir
			}

		}

		//

	} else {
		config.IPConfig = defaultConfig.IPConfig
		config.RecipeDir = defaultConfig.RecipeDir

	}

	//try flag recipeDir
	if *flagRecipeDir != defaultRecipeDir {
		infoLogger.Printf("trying to use recipeDir %s from flag", *flagRecipeDir)
		//MkdirAll returns nil on already exists or dir created, therefore errors are
		//serious
		mkErr := os.MkdirAll(*flagRecipeDir, 0644)
		//if error creating recipeDir, then attempt to use config set if not default
		if mkErr != nil {
			infoLogger.Printf("Unable to use %s as recipe directory, normally permissions or something out of my control", *flagRecipeDir)
		}
	} else if config.RecipeDir != defaultRecipeDir {
		infoLogger.Printf("trying to use recipeDir %s from config", config.RecipeDir)
		//MkdirAll returns nil on already exists or dir created, therefore errors are
		//serious
		mkErr := os.MkdirAll(config.RecipeDir, 0644)
		//if error creating recipeDir, then attempt to use default
		if mkErr != nil {
			infoLogger.Printf("Unable to use %s as recipe directory, normally permissions or something out of my control", config.RecipeDir)
		}
	} else {
		infoLogger.Printf("trying to use default recipeDir %s", defaultRecipeDir)
		//MkdirAll returns nil on already exists or dir created, therefore errors are
		//serious
		mkErr := os.MkdirAll(defaultRecipeDir, 0644)
		//if error creating recipeDir, then attempt to use default
		if mkErr != nil {
			infoLogger.Printf("Unable to use %s as recipe directory, normally permissions or something out of my control", defaultRecipeDir)
			//return err for fatal processing in main
			return mkErr
		}
	}
	return nil
}

//read recipe files into memory
func readRecipes(dirPath string) {
	readRecipe := func(path string, f os.FileInfo, err error) error {
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			//fmt.Println("Is directory: ", path)
			return nil
		}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		var recipe backend.Recipe

		err2 := json.Unmarshal(bytes, &recipe)
		if err2 != nil {
			return err
		}
		recipe.FileName = path
		recipes = append(recipes, recipe)
		return nil

	}
	err := filepath.Walk(dirPath, readRecipe)
	if err != nil {
		fatalLogger.Fatal("Could not walk directory", err)
	}
}

func saveRecipes(recipes []backend.Recipe) {
	for _, recipe := range recipes {
		tempPath := recipe.FileName
		bytes, err := json.MarshalIndent(recipe, "", "  ")
		if err != nil {
			infoLogger.Println("Error in saveRecipes(making json): ", err)
		}
		ioutil.WriteFile(tempPath, bytes, 0644)
		if err != nil {
			infoLogger.Println("Error in saveRecipes(writing file): ", err)
		}
	}

}

//check for and read ingredient database

//establish ncurses-like gui, look at termbox-go or gocui, mop-tracker
//or hecate for example programs

//add new recipe function

//addRecipe is used for command line stuff. uses os.stdin to progressively retrieve information
//returns a backend.recipe
func addRecipe() backend.Recipe {

	return backend.Recipe{}
}

func displaySingleRecipe(recipeName string) error {
	tempRecipe := backend.Recipe{}
	for _, recipe := range recipes {
		if recipe.Name == recipeName {
			tempRecipe = recipe
			break
		}
	}
	if tempRecipe.Name == "" {
		err := errors.New("Recipe not found")
		return err
	}
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
