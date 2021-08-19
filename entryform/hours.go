package entryform

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/jroimartin/gocui"
)

type HoursView struct {
	gui   *gocui.Gui
	hours chan int
}

func (hv *HoursView) handleEnter(hours int) {
	hv.gui.DeleteView("hours")
	hv.hours <- hours
}

func (hv *HoursView) editorFunc(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case mod == 0 && unicode.IsNumber(ch):
		v.EditWrite(ch)
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyEnter:
		hours, _ := strconv.Atoi(strings.TrimSpace(v.Buffer()))
		hv.handleEnter(hours)
		return
	}
}

func (hv *HoursView) GetHours(g *gocui.Gui) chan int {
	hv.hours = make(chan int)

	hv.gui = g
	maxX, _ := g.Size()

	if v, err := g.SetView("hours", 1, 2, maxX/2-2, 4); err != nil {
		v.Editable = true
		v.Title = "Hours"

		if _, err := g.SetCurrentView("hours"); err != nil {
			return nil
		}

		v.Editor = gocui.EditorFunc(hv.editorFunc)
	}

	return hv.hours
}
