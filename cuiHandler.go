package main

import (
	"database/sql"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	app           *tview.Application
	recipeList    *tview.List
	tagList       *tview.List
	recipePreview *tview.TextView
	//cmdPalette        *tview.TextView
	//recipeViewer      *tview.TextView
	mainFlexContainer *tview.Flex
	//recipeForm        *tview.Form
)

const (
	navigation = "Use arrow keys to navigate.\n" +
		"Left and right to switch between tags andrecipes,\n" +
		"and up and down arrows to select either a tag or recipe."
	commands = "Ctrl-C: Exit, "
)

// startCUI initializes CUI and starts it afterwards
func startCUI(db *sql.DB) error {

	// want 3 column layout

	// if possible, focus follows mouse, and can be shifted with keybindings

	// default layout is alphabetical list of recipe names on left,
	// commands listed in bottom half of center column,
	// top half blank
	// right column shows all tags

	// When recipe is selected in left column, name and description are
	// shown in top half of center column, with key commands shifted down
	// to bottom half

	// key command opens selected recipe

	// if recipe is opened, middle column shows
	// title, total time, instructions
	// left column shows ingredient list
	// right column is split, with current recipe tags on top
	// (selectable/searchable) and currently available commands on bottom.

	// when recipe is closed, view goes back to default

	// if tag is selected, recipes shown on left are filtered to only those
	// with the selected tag. Allow for multiple tag selection which will
	// do logical AND selection

	// search function opens modal dialog that does full text search of
	// titles/descriptions/comments/authors/sources (selectable)
	// and then populates a separate modal list with the results.
	// options are, open/refine search/escape

	// edit and create recipe use same form. Edit recipe is accessible when a
	// recipe is opened. Keybinds to edit recipe, and create new linked revision of recipe
	// new recipe is only available from main screen and opens a blank form

	// need interface for managing inventory

	// first create cui struct

	//var err error

	// then create application
	app = tview.NewApplication()

	// then create all interface widgets / primitives

	// lists are used for recipe and tag views.
	// hide secondary text, to just show names
	cmdPalette, err := initCmdPalette()
	if err != nil {
		return err
	}

	recipeList, err = initRecipeList(db)
	if err != nil {
		return err
	}
	tagList, err = initTagList(db)
	if err != nil {
		return err
	}
	recipePreview, err = initRecipePreview()
	if err != nil {
		return err
	}
	//recipeViewer, err = openRecipeViewer() //TODO: need to pass recipe into this
	//if err != nil {return err}

	//recipeEntryForm = initRecipeForm(false) // TODO: how to get values out of this
	// set up keybinds
	recipeList.SetInputCapture(recipeListKeybinds) //passing in function pointer

	tagList.SetInputCapture(tagListKeybinds) // passing in function pointer
	mainFlexContainer := tview.NewFlex()

	// set up initial view

	mainFlexContainer.AddItem(recipeList, 0, 1, true)
	// add in internal flexbox to get two row layout
	mainFlexContainer.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(recipePreview, 0, 2, false).
		AddItem(cmdPalette, 0, 1, false), 0, 2, false)
	mainFlexContainer.AddItem(tagList, 0, 1, false)

	err = app.SetRoot(mainFlexContainer, true).SetFocus(mainFlexContainer).Run()

	if err != nil {
		return err
	}

	return nil
}

func tagListKeybinds(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:
		item := tagList.GetCurrentItem()
		tagList.SetCurrentItem(item - 1)

		return nil
	case tcell.KeyDown:
		item := tagList.GetCurrentItem()
		tagList.SetCurrentItem(item + 1)

		return nil
	// only two places to go, so left or right doesn't matter
	case tcell.KeyLeft, tcell.KeyRight:
		app.SetFocus(recipeList)
		return nil

	case tcell.KeyRune:
		switch event.Rune() {
		case 'o', 'O':
			// open recipe
			return nil
		case 'e', 'E':
			//edit recipe
			return nil

		}

	}

	return event

}

func recipeListKeybinds(event *tcell.EventKey) *tcell.EventKey {

	switch event.Key() {
	case tcell.KeyUp:
		item := recipeList.GetCurrentItem()
		recipeList.SetCurrentItem(item - 1)

		return nil
	case tcell.KeyDown:
		item := recipeList.GetCurrentItem()
		recipeList.SetCurrentItem(item + 1)

		return nil
	// only two places to go, so left or right doesn't matter
	case tcell.KeyLeft, tcell.KeyRight:
		app.SetFocus(tagList)
		return nil

	case tcell.KeyRune:
		switch event.Rune() {
		case 'o', 'O':
			// open recipe
			return nil
		case 'e', 'E':
			//edit recipe
			return nil
		case 's', 'S':
			//find recipe
			return nil
		case 'n', 'N':
			//add recipe

			return nil

		}

	}

	return event

}

func initRecipeList(db *sql.DB) (*tview.List, error) {
	// may want to use treeView for recipes if we can have multiple versions of them
	recipeList = tview.NewList()
	recipeList.ShowSecondaryText(false).SetWrapAround(true)
	recipeList.SetBorder(true).SetTitle("Recipes")
	// always highlight selected list item
	recipeList.SetSelectedFocusOnly(false)

	recipes, err := GetRecipes(db)

	//TODO: need to sort recipe names

	for id, recipe := range recipes {
		recipeList.AddItem(recipe, strconv.Itoa(id), 0, nil)

	}

	return recipeList, err

}

func initTagList(db *sql.DB) (*tview.List, error) {

	tagList = tview.NewList()
	tagList.ShowSecondaryText(false).SetWrapAround(true)
	tagList.SetBorder(true).SetTitle("Tags")
	// always highlight selected list item
	tagList.SetSelectedFocusOnly(false)

	tagNames, err := GetTags(db)

	//TODO: need to sort tags

	for id, tag := range tagNames {

		tagList.AddItem(tag, srconv.Itoa(id), 0, nil)
	}

	return tagList, err
}

func initRecipePreview(recipeID int, db *sql.DB) (*tview.TextView, error) {

	recipePreview := tview.NewTextView()
	recipePreview.SetDynamicColors(true).SetRegions(false)
	recipePreview.SetChangedFunc(func() { app.Draw() })
	recipePreview.SetBorder(true).SetTitle("Selected Recipe")

	return recipePreview, nil
}

func initCmdPalette() (*tview.TextView, error) {

	cmdPalette := tview.NewTextView()
	cmdPalette.SetScrollable(true).SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorWhite)
	cmdPalette.SetWrap(true).SetWordWrap(true)
	cmdPalette.SetDynamicColors(true).SetRegions(false).SetChangedFunc(func() { app.Draw() })
	cmdPalette.SetText(navigation + "\n\n" + commands)

	return cmdPalette, nil
}

// openRecipeViewer opens a pretty printed view only view of a recipe
// suitable for actually making the recipe.
func openRecipeViewer(recipeID int) (*tview.TextView, error) {

	recipeViewer := tview.NewTextView()
	// TODO: embed regions and colors when recipes are read from database
	recipeViewer.SetDynamicColors(true).SetRegions(true)
	recipeViewer.SetBorder(true)
	// TODO: set title

	return recipeViewer, nil

}

// openRecipeForm allows for entering or editing recipes.
//
// edit parameter allows for the use of this function to add and edit recipes.
// Pass in false and -1 to add a new recipe, pass in true and the recipe id to edit
// the int return value will always return the recipe id that was either edited or created
func openRecipeForm(edit bool, recipeID int, db *sql.DB) (*tview.Form, int, error) {

	// setting up form for adding recipes
	recipeForm := tview.NewForm()
	// name, initial value, field width, validation function, changed function
	recipeForm.AddInputField("Name", "", 30, nil, nil)
	recipeForm.AddInputField("Description", "", 60, nil, nil)
	recipeForm.AddInputField("Comments", "", 300, nil, nil)
	recipeForm.AddInputField("Source", "", 30, nil, nil)
	recipeForm.AddInputField("Author", "", 30, nil, nil)
	recipeForm.AddInputField("Quantity", "", 4, nil, nil) //TODO: add validation function so it only allows numbers
	// label, options[], initialoption, selected function)
	qtyUnits := []string{"cookie", "serving", "Bar"} //TODO: get from database
	recipeForm.AddDropDown("Units", qtyUnits, 0, nil)
	//ingredientTable := tview.NewTable()
	//ingredientTable.SetFixed(3,2) // number of (rows,columns) visible at all times
	//ingredientTable.SetSeparator(tview.GraphicsVertBar)
	//numIngredients := IngredientCount()
	//recipeEntryForm.AddFormItem("Ingredients",)
	//recipeEntryForm.AddButton("Add Ingredient", nil)

	return recipeForm, -1, nil
}
