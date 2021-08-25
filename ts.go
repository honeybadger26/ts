package main

import (
	"log"
	"ts/entryform"

	"github.com/jroimartin/gocui"
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func initialiseViews(g *gocui.Gui) {
	NewHelp(g)
	NewLogger(g)

	entries := NewEntries(g)
	form := entryform.NewForm(g)
	info := NewInfo(g)

	form.HandleItemChange = info.UpdateInfo
	form.HandleEntrySubmitted = entries.RefreshEntries

	go func() {
		for {
			form.AddEntry()
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

	initialiseViews(g)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
