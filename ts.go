package main

import (
	"fmt"
	"log"

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

func printItemInfo(g *gocui.Gui, itemName string) {
	db := Database{}
	item := db.getItem(itemName)

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(INFO_VIEW)

		if err != nil {
			return err
		}

		v.Clear()

		if item == nil {
			return nil
		}

		fmt.Fprintf(v, "Name:           %s\n", item.Name)
		fmt.Fprintf(v, "Description:    %s\n", item.Description)
		if item.Size != "" {
			fmt.Fprintf(v, "Size:           %s\n", item.Size)
		}
		if item.TotalHours != -1 {
			fmt.Fprintf(v, "Total Hours:    %f\n", item.TotalHours)
		}

		return nil
	})
}

var HELP_TEXT = map[string]string{
	ITEM_VIEW: `<Ctrl-j> Select next item
<Ctrl-k> Select previous item
<Enter> Confirm selected item
<Alt-l> Next category
<Alt-h> Previous category`,
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

func printEntries(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(ENTRIES_VIEW)

		if err != nil {
			return err
		}

		v.Clear()

		db := Database{}
		entries := db.getEntries()
		cols, rows := v.Size()
		padding := cols/3 - 1

		for i, e := range entries {
			rowText := fmt.Sprintf("%-*s %-*s %*d", padding, e.Date, padding, e.Item, padding, e.Hours)

			if i%2 == 0 {
				rowText = fmt.Sprintf("\x1b[0;33m%s\x1b[0m", rowText)
			}

			fmt.Fprintln(v, rowText)
		}

		for len(v.BufferLines()) < rows {
			fmt.Fprintln(v, "")
		}

		pretext := "Total"
		fmt.Fprintf(v, "\x1b[0;32m%s %*d\x1b[0m", pretext, cols-len(pretext)-2, db.getTotalHours())

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

	initApp(g)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
