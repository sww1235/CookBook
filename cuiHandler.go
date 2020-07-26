package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	app               *tview.Application
	recipeList        *tview.List
	tagList           *tview.List
	recipePreview     *tview.TextView
	cmdPalette        *tview.Frame
	recipeViewer      *tview.TextView
	mainFlexContainer *tview.Flex
)

const (
	navigation = "Use arrow keys to navigate. Left and right to switch between tags and " +
		"recipes, and up and down arrows to select either a tag or recipe."
	commands = "Ctrl-C: Exit, "
)

func startCUI() error {

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
	//
	//

	// first create application
	app = tview.NewApplication()
	// then create all interface widgets / primitives

	// lists are used for recipe and tag views.
	// hide secondary text, to just show names

	// may want to use treeView for recipes if we can have multiple versions of them
	recipeList = tview.NewList()
	recipeList.ShowSecondaryText(false).SetWrapAround(true)
	recipeList.SetBorder(true).SetTitle("Recipes")
	// always highlight selected list item
	recipeList.SetSelectedFocusOnly(false)

	tagList = tview.NewList()
	tagList.ShowSecondaryText(false).SetWrapAround(true)
	tagList.SetBorder(true).SetTitle("Tags")
	// always highlight selected list item
	tagList.SetSelectedFocusOnly(false)

	recipePreview = tview.NewTextView()
	recipePreview.SetDynamicColors(true).SetRegions(false)
	recipePreview.SetChangedFunc(func() { app.Draw() })
	recipePreview.SetBorder(true).SetTitle("Selected Recipe")

	// using box as placeholder inside frame.
	// we can add text to frame but not box for some reason
	// frame needs primitive to go around
	// this comes from the presentation demo cover.go
	cmdPalette = tview.NewFrame(tview.NewBox())
	cmdPalette.SetBorders(0, 0, 0, 0, 0, 0)
	cmdPalette.AddText(navigation, true, tview.AlignCenter, tcell.ColorWhite)
	cmdPalette.AddText("", true, tview.AlignCenter, tcell.ColorWhite)

	recipeViewer = tview.NewTextView()
	// TODO: embed regions and colors when recipes are read from database
	recipeViewer.SetDynamicColors(true).SetRegions(true)
	recipeViewer.SetBorder(true)
	// TODO: set title

	// set up keybinds
	recipeList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

			}

		}

		return event

	})

	tagList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

			}

		}

		return event

	})
	mainFlexContainer = tview.NewFlex()

	// set up initial view

	mainFlexContainer.AddItem(recipeList, 0, 1, true)
	// add in internal flexbox to get two row layout
	mainFlexContainer.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(recipePreview, 0, 2, false).
		AddItem(cmdPalette, 0, 1, false), 0, 2, false)
	mainFlexContainer.AddItem(tagList, 0, 1, false)

	err := app.SetRoot(mainFlexContainer, true).SetFocus(mainFlexContainer).Run()

	if err != nil {
		fatalLogger.Panicln(err)
	}

	return nil
}
