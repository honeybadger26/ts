package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jroimartin/gocui"
)

/////////////////////////////////////////////////////////////

type Controller struct {
	gui *gocui.Gui

	logger      *Logger
	itemView    *ItemView
	hoursView   *HoursView
	messageView *MessageView
	entries     *Entries
}

// TODO: form should be it's own component
func (c *Controller) InitForm() {
	c.itemView = &ItemView{}
	c.hoursView = &HoursView{}
	c.messageView = &MessageView{}

	c.itemView.Init(c)
	c.hoursView.Init(c)
	c.messageView.Init(c)
}

func (c *Controller) Init(g *gocui.Gui) {
	c.gui = g
	c.logger = &Logger{}
	c.entries = &Entries{}

	c.logger.Init(c)
	c.entries.Init(c)

	c.InitForm()
}

func (c *Controller) NewEntry() {
	c.itemView.Destroy()
	c.hoursView.Destroy()
	c.messageView.Destroy()
	c.InitForm()
}

func (c *Controller) SubmitLog() {
	t := time.Now()
	line := fmt.Sprintf("%s, %s, %d\n", t.Format("02/01/2006"), c.itemView.Item, c.hoursView.Hours)

	f, err := os.OpenFile("data/savedlogs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()

	if _, err := f.WriteString(line); err != nil {
		log.Panicln(err)
	}

	c.NewEntry()
}

/////////////////////////////////////////////////////////////

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	c := &Controller{}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = false

	c.Init(g)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
