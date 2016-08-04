package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func startCUI() error {
	gui := gocui.NewGui()
	err := gui.Init()
	if err != nil {
		return err
	}
	defer gui.Close()

	gui.SetLayout(layout)

	if err := initKeybindings(gui); err != nil {
		return err
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func layout(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()

	mainView, mainErr := gui.SetView("main", 0, 0, maxX, maxY)
	if mainErr != nil {
		if mainErr != gocui.ErrUnknownView {
			return mainErr
		}
		fmt.Fprintln(mainView, "^C: Exit")
	}
	return nil
}

func initKeybindings(gui *gocui.Gui) error {
	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}
