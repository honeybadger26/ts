package main

import (
	"ts/entryform"

	"github.com/jroimartin/gocui"
)

type Controller struct {
	gui *gocui.Gui

	logger   *Logger
	form     *entryform.Form
	entries  *Entries
	helpView *HelpView
	info     *InfoComponent
}

func (c *Controller) Init(g *gocui.Gui) {
	c.gui = g
	c.logger = &Logger{}
	c.form = &entryform.Form{}
	c.entries = &Entries{}
	c.helpView = &HelpView{}

	c.info = NewInfo(g)
	c.logger.Init(c)
	c.form.Init(g, c.info.UpdateInfo, c.entries.RefreshEntries)
	c.entries.Init(c)
	c.helpView.Init(c)

	go func() {
		for {
			c.form.AddEntry()
		}
	}()
}
