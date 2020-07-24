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

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func layout(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()

	cmdView, cmdErr := gui.SetView("cmd", 0, maxY-2, maxX, maxY, 0)
	if cmdErr != nil {
		if cmdErr != gocui.ErrUnknownView {
			return cmdErr
		}
		fmt.Fprintln(cmdView, "^C: Exit")
	}
	mainView, mainErr := gui.SetView("main", 0, 0, maxX, maxY-2, 0)
	if mainErr != nil {
		if mainErr != gocui.ErrUnknownView {
			return mainErr
		}
		fmt.Fprintln(mainView, "this is a test")
	}
	return nil
}

func initKeybindings(gui *gocui.Gui) error {
	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := gui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return test3(g, v)
		}); err != nil {
		return err
	}

	return nil
}

func test2(gui *gocui.Gui, view *gocui.View) error {
	infoLogger.Println(view.Name())
	return test3(gui, view)
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
