package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type HelpView struct{}

func NewHelp(g *gocui.Gui) *HelpView {
	h := &HelpView{}
	maxX, maxY := g.Size()

	helpstr := "" +
		"<Ctrl-j> (when entering item) select next item \n" +
		"<Ctrl-k> (when entering item) select previous item \n" +
		"<Ctrl-c> to quit"

	if v, err := g.SetView("help", 0, maxY-2-3, (maxX/2)-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return nil
		}

		v.Wrap = true
		v.Frame = false

		fmt.Fprintf(v, helpstr)
	}

	return h
}
