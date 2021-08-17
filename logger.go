package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Logger struct {
	controller *Controller
}

func (l *Logger) Init(c *Controller) error {
	l.controller = c

	g := c.gui
	maxX, maxY := g.Size()

	if v, err := g.SetView("LOGGING", maxX/2, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Editable = false
		v.Frame = true
		v.Title = "LOG"
	}

	return nil
}

func (l *Logger) Log(text string) {
	l.controller.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("LOGGING")

		if err != nil {
			return err
		}

		fmt.Fprintln(v, text)
		return nil
	})
}
