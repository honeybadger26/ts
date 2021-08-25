package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type HelpView struct{}

func NewHelp(g *gocui.Gui) *HelpView {
	h := &HelpView{}
	maxX, maxY := g.Size()

	if v, err := g.SetView("help", 0, maxY-3, (maxX/2)-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return nil
		}

		v.Wrap = true
		v.Frame = false

		fmt.Fprintln(v, "<Ctrl-c> to quit")
	}

	return h
}
