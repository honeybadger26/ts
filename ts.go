package main

import (
	"log"
	"os"

	"github.com/jroimartin/gocui"

	"ts/editmode"
	"ts/viewmode"
)

func setKeyBindings(g *gocui.Gui) {
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = false

	setKeyBindings(g)

	if len(os.Args) == 1 {
		editmode.NewEditApp(g)
	} else if os.Args[1] == "v" {
		viewmode.NewViewApp(g, true)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
