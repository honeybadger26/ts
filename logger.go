package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Logger struct {
	gui *gocui.Gui
}

func NewLogger(g *gocui.Gui) *Logger {
	l := &Logger{}
	l.gui = g

	maxX, maxY := g.Size()

	if v, err := g.SetView("LOGGING", maxX/2, maxY-7, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return nil
		}
		v.Wrap = true
		v.Editable = false
		v.Frame = true
		v.Title = "LOG"
	}

	return l
}

func (l *Logger) Log(text string) {
	l.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("LOGGING")

		if err != nil {
			return err
		}

		fmt.Fprintln(v, text)
		return nil
	})
}
