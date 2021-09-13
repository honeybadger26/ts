package main

import (
	"log"
	"os"

	"github.com/jroimartin/gocui"
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
		NewApp(g)
	} else if os.Args[1] == "v" {
		NewViewApp(g)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		// 		log.Println("big ol errrrrr")
		// 		os.Exit(0)
		log.Panicln(err)
	}
}
