package main

//TODO: either split units out of ingredient.go and into their own file, or at
//least implement a way to add conversions between volume units (ie: cups to
//teespoons) and weight units (grams to lbs), as well as be able to add custom
//conversions between a specific weight unit and a specific volume unit for a
//specific ingredient

//TODO: figure out logging

import (
	"encoding/json"
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

var config Configuration

var recipes []backend.Recipe

var infoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	err := initialization()
	if err != nil {
		infoLogger.Panicln("Panic: Serious issue detected", err)
	}
	readRecipes(config.RecipeDir)

	writeConfig(config, path.Join(config.ConfigPath, "cookbook.cfg"))
}

//initialization sets up command line options, logging and config stuffs
func initialization() error {

	currUsr, usrErr := user.Current()
	if usrErr != nil {
		return usrErr
	}

	//TODO: Set flags to use temp variables so they can override config file

	//default paths. Will not be overridden
	configDir := path.Join(currUsr.HomeDir, ".config", "cookbook")
	recipeDir := path.Join(currUsr.HomeDir, "cookbook")

	config.ConfigPath = path.Join(configDir, "cookbook.cfg")
	config.RecipeDir = recipeDir

	defaultFileConfig, defaultCfgErr := readConfig(config.ConfigPath)

	//when this block is triggered, the config file is not in its default location
	if defaultCfgErr != nil {
		infoLogger.Println("Config file not found in default location, trying commandline flag")
		flag.StringVar(&config.ConfigPath, "c", config.ConfigPath, "Path to config file")
		flag.StringVar(&viewedRecipe, "v", "", "Recipe to view")
		flag.StringVar(&config.RecipeDir, "r", recipeDir, "Directory to store recipes in")
		flag.BoolVar(&addRecipeToggle, "n", false, "Add new recipe")
		flag.BoolVar(&httpServer, "H", false, "Use HTTP server instead of terminal")
		flag.StringVar(&config.IPConfig, "ip", "127.0.0.1", "IP to start HTTP server on")
		flag.Parse()

		mkErr := os.MkdirAll(config.ConfigPath, 0644)
		mkErr2 := os.MkdirAll(config.RecipeDir, 0644)

		//MkdirAll returns nil on already exists or dir created, therefore errors are
		//serious
		if mkErr != nil {
			return mkErr
		}
		if mkErr2 != nil {
			return mkErr2
		}

		readFileConfig, readCfgErr := readConfig(config.ConfigPath)
		if readCfgErr != nil {

		}
	}

	mkErr := os.MkdirAll(config.ConfigPath, 0644)
	mkErr2 := os.MkdirAll(config.RecipeDir, 0644)

	//MkdirAll returns nil on already exists or dir created, therefore errors are
	//serious
	if mkErr != nil {
		return mkErr
	}
	if mkErr2 != nil {
		return mkErr2
	}

	config.ConfigPath = defaultFileConfig.ConfigPath
	config.IPConfig = defaultFileConfig.IPConfig
	config.RecipeDir = defaultFileConfig.RecipeDir
	return nil
}

//read config file, either from -c flag or default in ~/.config

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
		infoLogger.Fatal("FATAL: Could not walk directory", err)
	}
}

func saveRecipes(recipes []backend.Recipe) {
	for _, recipe := range recipes {
		tempPath := recipe.FileName
		bytes, err := json.MarshalIndent(recipe, "", "  ")
		if err != nil {
			fmt.Println("Error in saveRecipes(making json): ", err)
		}
		ioutil.WriteFile(tempPath, bytes, 0644)
		if err != nil {
			fmt.Println("Error in saveRecipes(writing file): ", err)
		}
	}

}

//check for and read ingredient database

//establish ncurses-like gui, look at termbox-go or gocui, mop-tracker
//or hecate for example programs

//add new recipe function

func addRecipe() backend.Recipe {

	return backend.Recipe{}
}

//view recipe function

//format recipe to markdown

//Functions to shutdown program
//-write out all new and changed recipes to files
