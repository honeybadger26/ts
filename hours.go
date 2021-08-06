package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/jroimartin/gocui"
)

type HoursView struct {
	controller *Controller
	Hours      int
}

func (hv *HoursView) Init(c *Controller) error {
	hv.controller = c
	return nil
}

func (hv *HoursView) Destroy() {
	hv.controller.gui.DeleteView("hours")
}

func (hv *HoursView) handleEnter() {
	g := hv.controller.gui
	maxX, _ := g.Size()

	v, _ := g.SetView("hours", 0, 3, maxX/2-1, 5)
	v.Editable = false
	v.Clear()
	fmt.Fprintf(v, "\x1b[0;34m%d", hv.Hours)
	g.Cursor = false

	// TODO: this needs to move
	hv.controller.messageView.ShowSubmittedMessage()
	hv.controller.SubmitLog()
	hv.controller.entries.RefreshEntries()
}

func (hv *HoursView) editorFunc(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case mod == 0 && unicode.IsNumber(ch):
		v.EditWrite(ch)
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyEnter:
		hv.handleEnter()
		return
	}

	hours, err := strconv.Atoi(strings.TrimSpace(v.Buffer()))
	if err != nil {
		// handle error
	}
	hv.Hours = hours
}

func (hv *HoursView) Focus() {
	g := hv.controller.gui
	maxX, _ := g.Size()

	if v, err := g.SetView("hours", 0, 4, maxX/2-1, 6); err != nil {
		v.Editable = true
		v.Title = "Hours"

		if _, err := g.SetCurrentView("hours"); err != nil {
			return
		}

		v.Editor = gocui.EditorFunc(hv.editorFunc)
	}
}
