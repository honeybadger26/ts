package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"ts/entryform"

	"github.com/jroimartin/gocui"
)

type Controller struct {
	gui *gocui.Gui

	logger   *Logger
	form     *entryform.Form
	entries  *Entries
	helpView *HelpView
}

func (c *Controller) Init(g *gocui.Gui) {
	c.gui = g
	c.logger = &Logger{}
	c.form = &entryform.Form{}
	c.entries = &Entries{}
	c.helpView = &HelpView{}

	c.logger.Init(c)
	c.form.Init(g)
	c.entries.Init(c)
	c.helpView.Init(c)

	go func() {
		for {
			c.SubmitLog()
		}
	}()
}

func (c *Controller) SubmitLog() {
	t := time.Now()
	item, hours := c.form.AddEntry()

	line := fmt.Sprintf("%s, %s, %d\n", t.Format("02/01/2006"), item, hours)

	f, err := os.OpenFile("data/savedlogs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()

	if _, err := f.WriteString(line); err != nil {
		log.Panicln(err)
	}

	c.entries.RefreshEntries()
}
