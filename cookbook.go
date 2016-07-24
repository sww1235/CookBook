package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	backend "github.com/sww1235/recipe-database-backend"
)

var configPath string
var viewedRecipie string
var addRecipeToggle bool
var httpServer bool
var ipConfig string

var recipes []backend.Recipe

func main() {

	initCMD()
}

//read commandline options
func initCMD() {
	flag.StringVar(&configPath, "c", "~/.config/", "Path to config directory")
	flag.StringVar(&viewedRecipie, "v", "", "Recipe to view")
	flag.BoolVar(&addRecipeToggle, "n", false, "Add new recipe")
	flag.BoolVar(&httpServer, "H", false, "Start HTTP server on localhost")
	flag.StringVar(&ipConfig, "ip", "127.0.0.1", "IP to start HTTP server on")
	flag.Parse()
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
			fmt.Println("Is directory: ", path)
		}
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
