package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func logger(g *gocui.Gui, text string) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(LOGGER_VIEW)

		if err != nil {
			return err
		}

		fmt.Fprintln(v, text)

		_, vHeight := v.Size()
		if len(v.BufferLines()) > vHeight+1 {
			originx, originy := v.Origin()
			v.SetOrigin(originx, originy+1)
		}

		return nil
	})
}

type ItemData struct {
	Name        string
	Description string
	Size        string
	TotalHours  float32
}

func printItemInfo(g *gocui.Gui, itemName string) {
	file, err := os.Open("data/items.json")

	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	found := false
	var item ItemData

	decoder.Token()

	for decoder.More() {
		decoder.Decode(&item)
		if item.Name == itemName {
			found = true
			break
		}
	}

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(INFO_VIEW)

		if err != nil {
			return err
		}

		v.Clear()

		if !found {
			return nil
		}

		fmt.Fprintf(v, "Name:           %s\n", item.Name)
		fmt.Fprintf(v, "Description:    %s\n", item.Description)
		fmt.Fprintf(v, "Size:           %s\n", item.Size)
		fmt.Fprintf(v, "Total Hours:    %f\n", item.TotalHours)

		return nil
	})
}

var HELP_TEXT = map[string]string{
	ITEM_VIEW: `<Ctrl-j> Select next item
<Ctrl-k> Select previous item`,
	"APP": `<Ctrl-c> Quit`,
}

func printHelp(g *gocui.Gui, view string) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(HELP_VIEW)

		if err != nil {
			return err
		}

		v.Clear()

		if helpText, ok := HELP_TEXT[view]; ok {
			fmt.Fprintln(v, helpText)
		}

		fmt.Fprintf(v, HELP_TEXT["APP"])

		_, rows := v.Size()
		for len(v.BufferLines()) != rows {
			v.SetCursor(0, 0)
			v.EditNewLine()
		}

		return nil
	})
}

const SAVED_LOGS_FILE = "data/savedlogs.json"

type EntryData struct {
	Date  string
	Item  string
	Hours int
}

func printEntries(g *gocui.Gui) {
	file, _ := os.Open(SAVED_LOGS_FILE)
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.Token()
	var entries []EntryData
	var entry EntryData

	for decoder.More() {
		decoder.Decode(&entry)
		entries = append(entries, entry)
	}

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(ENTRIES_VIEW)

		if err != nil {
			return err
		}

		v.Clear()
		cols, _ := v.Size()
		padding := cols/3 - 1

		for i, e := range entries {
			rowText := fmt.Sprintf("%-*s %-*s %*d", padding, e.Date, padding, e.Item, padding, e.Hours)

			if i%2 == 0 {
				rowText = fmt.Sprintf("\x1b[0;33m%s\x1b[0m", rowText)
			}

			fmt.Fprintln(v, rowText)
		}

		return nil
	})
}

func initApp(g *gocui.Gui) {
	maxX, maxY := g.Size()

	for _, name := range MAIN_VIEWS {
		p := VIEW_PROPS[name]
		x0 := int(p.x0 * float32(maxX))
		y0 := int(p.y0 * float32(maxY))
		x1 := int(p.x1*float32(maxX)) - 1
		y1 := int(p.y1*float32(maxY)) - 1

		if v, err := g.SetView(name, x0, y0, x1, y1); err != nil {
			if err != gocui.ErrUnknownView {
				log.Panicln(err)
			}
			v.Title = p.title
			v.Wrap = true
			v.Editable = p.editable
			v.Frame = p.frame
		}
	}

	printHelp(g, "")
	printEntries(g)
	printItemInfo(g, "")

	go func() {
		for {
			addNewEntry(g)
			printEntries(g)
		}
	}()
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = false

	// initialiseViews(g)
	initApp(g)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
