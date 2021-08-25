package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

type Entries struct {
	gui *gocui.Gui
}

func NewEntries(g *gocui.Gui) *Entries {
	e := &Entries{}
	e.gui = g

	maxX, maxY := g.Size()

	if v, err := g.SetView("entries", maxX/2, 0, maxX-1, (maxY/2)-1); err != nil {
		if err != gocui.ErrUnknownView {
			return nil
		}
		v.Wrap = true
		v.Editable = false
		v.Frame = true
		v.Title = "Entries for today"
	}

	e.RefreshEntries()
	return e
}

func (e *Entries) RefreshEntries() {
	e.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("entries")

		if err != nil {
			return err
		}

		buf, err := os.Open("data/savedlogs")
		if err != nil {
			// fix this
			// c.logger.Log("Could not open savedlogs file. This will created")
			return nil
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
