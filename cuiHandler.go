package main

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

func startCUI() error {
	gui, err := gocui.NewGui(gocui.Output256, true)
	if err != nil {
		return err
	}
	defer gui.Close()

	gui.SetManagerFunc(layout)

	if err := initKeybindings(gui); err != nil {
		return err
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func layout(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()

	//cmd view prints available commands for each view
	cmdView, cmdErr := gui.SetView("cmd", 0, maxY-2, maxX, maxY, 0)
	if cmdErr != nil {
		if !gocui.IsUnknownView(cmdErr) {
			return cmdErr
		}
		fmt.Fprintln(cmdView, "^C: Exit")
	}

	// main view shows usage instructions and main keyboard commands
	mainView, mainErr := gui.SetView("main", 0, 0, maxX, maxY-2, 0)
	if mainErr != nil {
		if !gocui.IsUnknownView(mainErr) {
			return mainErr
		}
		fmt.Fprintln(mainView, "this is a test")
	}

	// recipe view displays individual recipe

	//recipeView, recipeErr := gui.SetView("recipe",

	return nil
}

func initKeybindings(gui *gocui.Gui) error {
	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := gui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, test3); err != nil {
		return err
	}

	return nil
}

func test3(gui *gocui.Gui, view *gocui.View) error {
	name := view.Name()
	infoLogger.Println(name)
	x0, y0, x1, y1, err := gui.ViewPosition(name)
	if err != nil {
		return err
	}
	fmt.Fprintln(view, "this is test2 firing")
	if _, err := gui.SetView(name, x0, y0, x1, y1, 0); err != nil {
		return err
	}

	return nil
}
