package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jroimartin/gocui"
)

type EntryForm struct {
	app *App

	// User input
	category ItemCategory
	item     string
	hours    int

	items         []Item
	filteredItems []Item
	selectedIndex int
	entry         Entry
}

func NewEntryForm(app *App) *EntryForm {
	ef := &EntryForm{}

	ef.app = app
	ef.category = All
	ef.item = ""
	ef.hours = 0

	ef.items = app.db.getItems(All)
	ef.filterItems()
	return ef
}

func (ef *EntryForm) changeItem(item string) {
	if item == ef.item {
		return
	}
	ef.item = item
	ef.filterItems()
}

func (ef *EntryForm) changeCategory(forward bool) {
	if forward {
		ef.category = (ef.category + 1) % ITEMTYPE_COUNT
	} else {
		ef.category = (ef.category + ITEMTYPE_COUNT - 1) % ITEMTYPE_COUNT
	}

	ef.items = ef.app.db.getItems(ef.category)
	ef.filterItems()
}

func (ef *EntryForm) filterItems() {
	ef.filteredItems = []Item{}

	for _, item := range ef.items {
		regexstr := `(?i)` + ef.item
		match, err := regexp.MatchString(regexstr, item.Name)
		if err != nil {
			// handle error
		}
		if match {
			ef.filteredItems = append(ef.filteredItems, item)
		}
	}

	if ef.selectedIndex = -1; len(ef.filteredItems) != 0 {
		ef.selectedIndex = 0
	}

	ef.updateItemView()
}

func (ef *EntryForm) updateItemView() {
	g := ef.app.gui

	if ef.selectedIndex != -1 {
		ef.app.changeItem(ef.filteredItems[ef.selectedIndex].Name)
	} else {
		ef.app.changeItem("")
	}

	fv, _ := g.View(FORM_VIEW)
	p := VIEW_PROPS[FORM_VIEW]
	maxX, maxY := g.Size()
	x0 := int(p.x0*float32(maxX)) + 1
	y0 := int(p.y0*float32(maxY)) + 1 + len(fv.BufferLines())
	x1 := int(p.x1*float32(maxX)) - 1 - 1
	y1 := int(p.y1*float32(maxY)) - 1 - 1

	iv, err := g.SetView(ITEM_VIEW, x0, y0, x1, y1)

	if err == nil {
		iv.Clear()
	} else if err != gocui.ErrUnknownView {
		return
	}

	iv.Wrap = true
	iv.Editable = VIEW_PROPS[ITEM_VIEW].editable
	iv.Frame = VIEW_PROPS[ITEM_VIEW].frame
	cols, rows := iv.Size()

	iv.Clear()

	if len(ef.filteredItems) == 0 {
		fmt.Fprintf(iv, "\x1b[0;31mNo results\x1b[0m\n")
	} else {
		for i, item := range ef.filteredItems {
			if i == ef.selectedIndex {
				fmt.Fprintf(iv, "\x1b[0;34m> %s\x1b[0m\n", item.Name)
			} else {
				fmt.Fprintln(iv, item.Name)
			}
		}
	}

	for len(iv.BufferLines()) < rows {
		fmt.Fprintln(iv, "")
	}

	fmt.Fprintf(iv, "%*s", cols, ef.category)
}

func (ef *EntryForm) changeSelectedIndex(forward bool) {
	indexBefore := ef.selectedIndex

	if forward && (indexBefore < len(ef.filteredItems)-1) {
		ef.selectedIndex++
	} else if !forward && (indexBefore > 0) {
		ef.selectedIndex--
	}

	if ef.selectedIndex != indexBefore {
		ef.updateItemView()
	}
}

func (ef *EntryForm) getItem() {
	done := make(chan bool)

	v, _ := ef.app.gui.View(FORM_VIEW)
	buffer := v.BufferLines()
	cX := len(buffer[len(buffer)-1])
	cY := len(buffer) - 1

	v.SetCursor(cX, cY)
	ef.updateItemView()

	v.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		switch {
		case key == gocui.KeyArrowDown || key == gocui.KeyCtrlJ:
			ef.changeSelectedIndex(true)
			return
		case key == gocui.KeyArrowUp || key == gocui.KeyCtrlK:
			ef.changeSelectedIndex(false)
			return
		case key == gocui.KeyEnter:
			if ef.app.item != "" {
				ef.app.gui.DeleteView(ITEM_VIEW)
				v.SetCursor(cX, cY)
				for range ef.filteredItems[ef.selectedIndex].Name {
					v.EditDelete(false)
				}
				done <- true
			}
			return
		case key == gocui.KeyBackspace || key == gocui.KeyBackspace2 || key == gocui.KeyArrowLeft:
			if newCursorX, newCursorY := v.Cursor(); newCursorX == cX && newCursorY == cY {
				return
			}
		case mod == gocui.ModAlt && ch == 'l':
			ef.changeCategory(true)
			return
		case mod == gocui.ModAlt && ch == 'h':
			ef.changeCategory(false)
			return
		}
		gocui.DefaultEditor.Edit(v, key, ch, mod)

		buf := v.BufferLines()
		line := buf[len(buf)-1]
		ef.changeItem(strings.TrimSpace(line[cX:]))
	})

	<-done
}

func (ef *EntryForm) getHours() int {
	hours := make(chan int)

	// put in function?
	v, _ := ef.app.gui.View(FORM_VIEW)
	buffer := v.BufferLines()
	cX := len(buffer[len(buffer)-1])
	cY := len(buffer) - 1

	v.SetCursor(cX, cY)

	v.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		switch {
		case ch != 0 && mod == 0 && !unicode.IsNumber(ch):
			return
		case key == gocui.KeyEnter:
			buf := v.BufferLines()
			line := buf[len(buf)-1]
			hoursStr := strings.TrimSpace(line[cX:])
			hoursInt, _ := strconv.Atoi(hoursStr)
			v.SetCursor(cX, cY)
			for range hoursStr {
				v.EditDelete(false)
			}
			hours <- hoursInt
			return
		case key == gocui.KeyBackspace || key == gocui.KeyBackspace2 || key == gocui.KeyArrowLeft:
			if newCursorX, newCursorY := v.Cursor(); newCursorX == cX && newCursorY == cY {
				return
			}
		}
		gocui.DefaultEditor.Edit(v, key, ch, mod)
	})

	return <-hours
}

func (ef *EntryForm) SetDate(date time.Time) {
	newDate := date.Format("02/01/2006")
	v, _ := ef.app.gui.View(FORM_VIEW)
	buffer := v.BufferLines()
	cX := len(buffer[len(buffer)-1])
	cY := len(buffer) - 1

	dateLineEnd := len(buffer[0])
	v.SetCursor(dateLineEnd, 0)
	for i := 0; i < dateLineEnd; i++ {
		v.EditDelete(true)
	}

	lineStr := "Date: " + newDate
	runeArr := []rune(lineStr)
	for _, ch := range runeArr {
		v.EditWrite(ch)
	}

	v.SetCursor(cX, cY)
	ef.entry.Date = newDate
}

func (ef *EntryForm) GetEntry() Entry {
	v, _ := ef.app.gui.View(FORM_VIEW)
	ef.app.gui.SetCurrentView(FORM_VIEW)
	v.Clear()

	e := &ef.entry

	date := ef.app.date.Format("02/01/2006")
	e.Date = date
	fmt.Fprintf(v, "Date: %s\n", date)

	v.Editable = true

	fmt.Fprintf(v, "Item: ")
	ef.getItem()
	item := ef.app.item
	e.Item = item
	fmt.Fprintln(v, item)

	fmt.Fprintf(v, "Hours: ")
	hours := ef.getHours()
	e.Hours = hours
	fmt.Fprintln(v, hours)

	v.Editable = false
	return *e
}
