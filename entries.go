package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

type Entries struct {
	controller *Controller
}

func (e *Entries) Init(c *Controller) error {
	e.controller = c

	g := c.gui
	maxX, maxY := g.Size()

	if v, err := g.SetView("entries", maxX/2, 0, maxX-1, (maxY/2)-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Editable = false
		v.Frame = true
		v.Title = "Entries for today"
	}

	e.RefreshEntries()
	return nil
}

func (e *Entries) RefreshEntries() {
	e.controller.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("entries")

		if err != nil {
			return err
		}

		buf, err := os.Open("data/savedlogs")
		if err != nil {
			log.Panicln(err)
		}

		defer func() {
			if err = buf.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		snl := bufio.NewScanner(buf)

		v.Clear()
		for snl.Scan() {
			fmt.Fprintln(v, snl.Text())
		}

		if err := snl.Err(); err != nil {
			log.Fatal(err)
		}

		return nil
	})
}
