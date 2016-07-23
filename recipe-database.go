package main

import "flag"

//read commandline options
func main() {
	flag.String("c", "~/.config/", "Path to config directory")
	flag.String("v", "", "Recipe to view")
	flag.String("n", "", "Add new recipe")
	flag.Bool("h", false, "Start HTTP server on localhost")
	flag.String("ip", "127.0.0.1", "IP to start HTTP server on")
}

//read config file, either from -c flag or default in ~/.config

//read recipe files into memory

//check for and read ingredient database

//establish ncurses-like gui, look at termbox-go or gocui, mop-tracker
//or hecate for example programs

//add new recipe function

//view recipe function

//format recipe to markdown

//Functions to shutdown program
//-write out all new and changed recipes to files
