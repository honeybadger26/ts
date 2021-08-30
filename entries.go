package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jroimartin/gocui"
)

const (
	SAVED_LOGS_FILE = "data/savedlogs.json"
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
		v.Title = "Entries"
	}

	e.RefreshEntries()
	return e
}

// repeated code - put this in own file
type EntryData struct {
	Date  string
	Item  string
	Hours int
}

func (e *Entries) RefreshEntries() {
	e.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("entries")
		v.Clear()
		cols, _ := v.Size()
		padding := cols/3 - 1

		if err != nil {
			return err
		}

		file, _ := os.Open(SAVED_LOGS_FILE)
		defer file.Close()

		decoder := json.NewDecoder(file)
		decoder.Token()
		var data EntryData

		count := 0
		for decoder.More() {
			decoder.Decode(&data)

			if count%2 == 0 {
				fmt.Fprintf(v, "%-*s %-*s %*d\n", padding, data.Date, padding, data.Item, padding, data.Hours)
			} else {
				fmt.Fprintf(v, "\x1b[0;33m%-*s %-*s %*d\x1b[0m\n", padding, data.Date, padding, data.Item, padding, data.Hours)
			}
			count++
		}

		return nil
	})
}
