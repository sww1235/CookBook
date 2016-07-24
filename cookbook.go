package main

//TODO: either split units out of ingredient.go and into their own file, or at
//least implement a way to add conversions between volume units (ie: cups to
//teespoons) and weight units (grams to lbs), as well as be able to add custom
//conversions between a specific weight unit and a specific volume unit for a
//specific ingredient

//TODO: figure out logging

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"

	backend "github.com/sww1235/recipe-database-backend"
)

var viewedRecipe string
var addRecipeToggle bool
var httpServer bool

var config Configuration

var recipes []backend.Recipe

func main() {
	flag.Parse()
	readConfig(path.Join(config.ConfigPath, "cookbook.cfg"))

	readRecipes(config.RecipePath)
}

//read commandline options
func init() {
	currUsr, usrErr := user.Current()
	if usrErr != nil {
		fmt.Println(usrErr)
	}
	configDir := path.Join(currUsr.HomeDir, ".config", "cookbook")
	recipeDir := path.Join(currUsr.HomeDir, "cookbook")

	//fmt.Println(usr.HomeDir)
	flag.StringVar(&config.ConfigPath, "c", configDir, "Path to config directory")
	flag.StringVar(&viewedRecipe, "v", "", "Recipe to view")
	flag.StringVar(&config.RecipePath, "r", recipeDir, "Directory to store recipes in")
	flag.BoolVar(&addRecipeToggle, "n", false, "Add new recipe")
	flag.BoolVar(&httpServer, "H", false, "Start HTTP server on localhost")
	flag.StringVar(&config.IPConfig, "ip", "127.0.0.1", "IP to start HTTP server on")
	flag.Parse()
	mkErr := os.MkdirAll(configDir, 740)
	//fmt.Println(configDir)
	if mkErr == nil {
		//fmt.Println("directory created or already exists")
	}
	if mkErr != nil {
		fmt.Println(mkErr)
	}
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
		fmt.Println("Is file: ", path)
		return nil

	}
	err := filepath.Walk(dirPath, readRecipe)
	if err != nil {
		fmt.Println("file does not exist")
	}

}

//check for and read ingredient database

//establish ncurses-like gui, look at termbox-go or gocui, mop-tracker
//or hecate for example programs

//add new recipe function

//view recipe function

//format recipe to markdown

//Functions to shutdown program
//-write out all new and changed recipes to files
