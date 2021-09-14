package main

import (
	"log"
	"os"

	"github.com/jroimartin/gocui"

	"ts/editmode"
	"ts/viewmode"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = false

	if len(os.Args) == 1 {
		editmode.NewApp(g)
	} else if os.Args[1] == "v" {
		viewmode.NewViewApp(g)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
