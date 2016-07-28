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
  cfgHandler "github.com/sww1235/go-config-handler"
)

//Configuration stores the configuration that is read in and out from a file
type Configuration struct {
	ConfigPath string `json:"configpath"`
	IPConfig   string
	RecipePath string
}

var viewedRecipe string
var addRecipeToggle bool
var httpServer bool


var config Configuration

var recipes []backend.Recipe

var infoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	initialization()

	readRecipes(config.RecipePath)

	cfgHandler.writeConfig(config, path.Join(config.ConfigPath, "cookbook.cfg"))
}

//initialization sets up command line options, logging and config stuffs
func initialization() error {

	currUsr, usrErr := user.Current()
	if usrErr != nil {
		return usrErr
	}

	//TODO: attempt to read config file, if it does not exist or can't be read,
	//then attempt to use flags, which then assign default values if not used.

	//default paths. Will not be overridden
	configDir := path.Join(currUsr.HomeDir, ".config", "cookbook")
	recipeDir := path.Join(currUsr.HomeDir, "cookbook")

	config.ConfigPath = configDir
	config.RecipePath = recipeDir

	//Have to change in program
	//flag.StringVar(&config.ConfigPath, "c", configDir, "Path to config directory")

	defaultFileConfig, defaultCfgErr := readConfig(path.Join(config.ConfigPath, "cookbook.cfg"))

	//when this block is triggered, the config file is not in its default location
	if defaultCfgErr != nil {
		infoLogger.Println("Config file not found in default location, trying")
		flag.StringVar(&viewedRecipe, "v", "", "Recipe to view")
		//Have to change in program
		//flag.StringVar(&config.RecipePath, "r", recipeDir, "Directory to store recipes in")
		flag.BoolVar(&addRecipeToggle, "n", false, "Add new recipe")
		flag.BoolVar(&httpServer, "H", false, "Start HTTP server on localhost")
		flag.StringVar(&config.IPConfig, "ip", "127.0.0.1", "IP to start HTTP server on")
		flag.Parse()

		readFileConfig, readCfgErr := readConfig(path.Join(config.ConfigPath, "cookbook.cfg"))
    if readCfgErr
	}

	mkErr := os.MkdirAll(config.ConfigPath, 0644)
	mkErr2 := os.MkdirAll(config.RecipePath, 0644)

	//fmt.Println(configDir)

	if mkErr != nil {
		return mkErr
	}
	if mkErr2 != nil {
		return mkErr2
	}
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
		fmt.Println("Error in ReadRecipes: ", err)
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
