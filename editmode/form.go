package editmode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jroimartin/gocui"

	"ts/database"
	"ts/viewmanager"
)

const (
	DISPLAY_DATE_FORMAT = "Mon 02 / Jan 01 / 2006"
	DATE_FORMAT         = "02/01/2006"
)

type EditMode int

const (
	emHourly EditMode = iota
	emDateRange
)

type EntryForm struct {
	app *App

	// User input
	category database.ItemCategory
	date     time.Time
	editMode EditMode
	item     string
	hours    int

	items         []database.Item
	filteredItems []database.Item
	selectedIndex int
}

type DateRange struct {
	from time.Time
	to   time.Time
}

// make a refreshView method. pretty much updateItemView but also:
// - redraws entire view (including item view)
// - sets cursor to end of buffer

func NewEntryForm(app *App) (ef *EntryForm) {
	ef = &EntryForm{}

	ef.app = app
	ef.editMode = emHourly
	ef.category = database.ICAll
	ef.item = ""
	ef.hours = 0

	ef.items = app.db.GetItems(ef.category)
	ef.filterItems()
	return
}

func (ef *EntryForm) changeItem(item string) {
	if item == ef.item {
		return
	}
	ef.item = item
	ef.filterItems()
}

func (ef *EntryForm) changeNextCategory() {
	ef.category = ef.category.GetNext()
	ef.items = ef.app.db.GetItems(ef.category)
	ef.filterItems()
}

func (ef *EntryForm) filterItems() {
	ef.filteredItems = []database.Item{}

	for _, item := range ef.items {
		match, _ := regexp.MatchString(`(?i)`+ef.item, item.Name+item.Description)
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
	if ef.selectedIndex != -1 {
		ef.app.changeItem(ef.filteredItems[ef.selectedIndex].Name)
	} else {
		ef.app.changeItem("")
	}

	g := ef.app.gui
	fv, _ := g.View(FORM_VIEW)
	p := VIEW_PROPS[FORM_VIEW]
	maxX, maxY := g.Size()
	// shouldn't have to do all this work to get points
	// could store the form views points on ef?
	x0 := int(p.X0*float32(maxX)) + 1
	y0 := int(p.Y0*float32(maxY)) + 1 + len(fv.BufferLines())
	x1 := int(p.X1*float32(maxX)) - 1 - 1
	y1 := int(p.Y1*float32(maxY)) - 1 - 1

	iv, err := g.SetView(ITEM_VIEW, x0, y0, x1, y1)

	if err == nil {
		iv.Clear()
	} else if err != gocui.ErrUnknownView {
		return
	}

	iv.Wrap = true
	iv.Editable = VIEW_PROPS[ITEM_VIEW].Editable
	iv.Frame = VIEW_PROPS[ITEM_VIEW].Frame

	iv.Clear()

	categoryText := ""
	for c := database.ICRecent; ; c = c.GetNext() {
		if categoryText != "" {
			categoryText += " - "
		}
		if c == ef.category {
			categoryText += fmt.Sprintf("\x1b[0;33m%s\x1b[0m", c)
		} else {
			categoryText += c.String()
		}

		if c.GetNext() == database.ICRecent {
			break
		}
	}

	cols, _ := iv.Size()
	// used to remove the padding that is added because of the \x1b stuff
	escapeCount := len("\x1b[0;33m\x1b[0m")
	fmt.Fprintf(iv, "%*s\n", cols-1+escapeCount, categoryText)

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
	cX, cY := viewmanager.GetEndPos(v)

	v.SetCursor(cX, cY)
	ef.updateItemView()

	v.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		switch {
		case key == gocui.KeyArrowDown: // most of these should be done with setKeyBindings
			ef.changeSelectedIndex(true)
			return
		case key == gocui.KeyArrowUp:
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
		case key == gocui.KeyCtrlD:
			if (ef.app.item != "") && (strings.Contains(ef.app.item, "Leave") || strings.Contains(ef.app.item, "Holiday")) {
				ef.app.gui.DeleteView(ITEM_VIEW)
				v.SetCursor(cX, cY)
				for range ef.filteredItems[ef.selectedIndex].Name {
					v.EditDelete(false)
				}
				ef.editMode = emDateRange
				done <- true
			}
			return
		case key == gocui.KeyBackspace || key == gocui.KeyBackspace2 || key == gocui.KeyArrowLeft:
			if newCursorX, newCursorY := v.Cursor(); newCursorX == cX && newCursorY == cY {
				return
			}
		case key == gocui.KeyTab:
			ef.changeNextCategory()
			return
		}
		gocui.DefaultEditor.Edit(v, key, ch, mod)

		item := viewmanager.GetEndString(v, cX)
		ef.changeItem(item)
	})

	<-done
}

func (ef *EntryForm) getHours() int {
	hours := make(chan int)

	// put in function?
	v, _ := ef.app.gui.View(FORM_VIEW)
	cX, cY := viewmanager.GetEndPos(v)

	v.SetCursor(cX, cY)

	v.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		switch {
		case ch != 0 && mod == 0 && !unicode.IsNumber(ch):
			return
		case key == gocui.KeyEnter:
			hoursStr := viewmanager.GetEndString(v, cX)
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

func (ef *EntryForm) GetInputDate() time.Time {
	dateChan := make(chan time.Time)

	v, _ := ef.app.gui.View(FORM_VIEW)
	cX, cY := viewmanager.GetEndPos(v)

	v.SetCursor(cX, cY)
	v.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		switch {
		case ch != 0 && mod == 0 && ch != 47 && !unicode.IsNumber(ch):
			return
		case key == gocui.KeyEnter:
			dateStr := viewmanager.GetEndString(v, cX)
			date, _ := time.Parse(DATE_FORMAT, dateStr)
			v.SetCursor(cX, cY)
			for range dateStr {
				v.EditDelete(false)
			}
			dateChan <- date
			return
		case key == gocui.KeyBackspace || key == gocui.KeyBackspace2 || key == gocui.KeyArrowLeft:
			if newCursorX, newCursorY := v.Cursor(); newCursorX == cX && newCursorY == cY {
				return
			}
		}
		gocui.DefaultEditor.Edit(v, key, ch, mod)
	})

	return <-dateChan
}

func (ef *EntryForm) GetDateRange() (dr DateRange) {
	v, _ := ef.app.gui.View(FORM_VIEW)

	fmt.Fprintf(v, "Start Date: ")
	startDate := ef.GetInputDate()
	fmt.Fprintln(v, startDate.Format(DATE_FORMAT))

	fmt.Fprintf(v, "End Date: ")
	endDate := ef.GetInputDate()
	fmt.Fprintln(v, endDate.Format(DATE_FORMAT))

	dr = DateRange{startDate, endDate}
	return
}

func (ef *EntryForm) SetDate(date time.Time) {
	v, _ := ef.app.gui.View(FORM_VIEW)
	buffer := v.BufferLines()
	cX, cY := viewmanager.GetEndPos(v)

	// can use v.Write('\r') to clear line?
	dateLineEnd := len(buffer[0])
	v.SetCursor(dateLineEnd, 0)
	for i := 0; i < dateLineEnd; i++ {
		v.EditDelete(true)
	}

	// this has problems when trying to add colors. fix this (use refreshView solution)
	diffWeekMsg := ""
	yearNow, weekNow := time.Now().ISOWeek()
	y, w := date.ISOWeek()
	if w < weekNow || y < yearNow {
		diffWeekMsg = "   >>> PAST WEEK <<<"
	} else if w > weekNow || y > yearNow {
		diffWeekMsg = "   >>> FUTURE WEEK <<<"
	}

	lineStr := "Date: " + date.Format(DISPLAY_DATE_FORMAT) + diffWeekMsg
	runeArr := []rune(lineStr)
	for _, ch := range runeArr {
		v.EditWrite(ch)
	}

	v.SetCursor(cX, cY)
	ef.date = date
}

// AMS - Refactor so that can be used for date ranges.
func (ef *EntryForm) GetEntries() (entrySlice []database.Entry) {
	v, _ := ef.app.gui.View(FORM_VIEW)
	ef.app.gui.SetCurrentView(FORM_VIEW)
	v.Clear()

	// Set the date of the form based on the app, then display the date in User Friendly format

	fmt.Fprintf(v, "Date: %s\n", ef.app.date.Format(DISPLAY_DATE_FORMAT))

	v.Editable = true

	// Get user input for item selection
	fmt.Fprintf(v, "Item: ")
	ef.getItem()
	fmt.Fprintln(v, ef.app.item)
	ef.item = ef.app.item

	// Get user input for entry date and hours
	var dateRange DateRange
	if ef.editMode == emHourly {
		ef.date = ef.app.date
		fmt.Fprintf(v, "Hours: ")
		ef.hours = ef.getHours()
		// Or get user input for date range
	} else if ef.editMode == emDateRange {
		v.Clear()
		fmt.Fprintf(v, "Item: ")
		fmt.Fprintln(v, ef.app.item)
		fmt.Fprintln(v, "Logging for date range...")
		dateRange = ef.GetDateRange()
	}

	v.Editable = false

	// Preparing entry(s) based on editMode and user input
	var entry database.Entry

	entry.Item = ef.item
	if ef.editMode == emHourly {
		entry.Date = ef.date
		entry.Hours = ef.hours
		entrySlice = append(entrySlice, entry)
	} else {
		for d := dateRange.from; d.After(dateRange.to) == false; d = d.AddDate(0, 0, 1) {
			if (int(d.Weekday()) == SATURDAY) || (int(d.Weekday()) == SUNDAY) {
				continue
			}

			entry.Date = d
			entry.Hours = FULL_DAY
			entrySlice = append(entrySlice, entry)
		}
	}

	return entrySlice
}
