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
)

const (
	DISPLAY_DATE_FORMAT = "Mon 02 / Jan 01 / 2006"
	DATE_FORMAT = "02/01/2006"
)

type EntryForm struct {
	app *App

	// User input
	category database.ItemCategory
	item	 string
	hours    int

	items         []database.Item
	filteredItems []database.Item
	selectedIndex int
	entry         database.Entry
}

type DateRange struct {
	from time.Time
	to time.Time
}

// make a refreshView method. pretty much updateItemView but also:
// - saves current cursor pos
// - updates the form view
// - restore cursor pos
// OR
// - redraws entire view (including item view)
// - sets cursor to end of buffer

func NewEntryForm(app *App) (ef *EntryForm) {
	ef = &EntryForm{}

	ef.app = app
	// ams - add a check to show ICRecent if there are recent items ?
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

	iv.Title = fmt.Sprintf("%s", ef.category)
	iv.Wrap = true
	iv.Editable = VIEW_PROPS[ITEM_VIEW].editable
	iv.Frame = VIEW_PROPS[ITEM_VIEW].frame

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
		case key == gocui.KeyBackspace || key == gocui.KeyBackspace2 || key == gocui.KeyArrowLeft:
			if newCursorX, newCursorY := v.Cursor(); newCursorX == cX && newCursorY == cY {
				return
			}
		case key == gocui.KeyTab:
			ef.changeNextCategory()
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

func (ef *EntryForm) GetInputDate() time.Time {
	dateChan := make(chan time.Time)

	v, _ := ef.app.gui.View(FORM_VIEW)
	buffer := v.BufferLines()
	cX := len(buffer[len(buffer)-1])
	cY := len(buffer) - 1

	v.SetCursor(cX, cY)
	v.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		switch {
		case ch != 0 && mod == 0 && ch != 47 && !unicode.IsNumber(ch):
			return
		case key == gocui.KeyEnter:
			buf := v.BufferLines()
			line := buf[len(buf)-1]
			dateStr := strings.TrimSpace(line[cX:])
			date, _ := time.Parse(DATE_FORMAT, dateStr)
			v.SetCursor(cX, cY)
			for range dateStr {
				v.EditDelete(false)
			}
			dateChan <-date
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
	// should make the date format a constant somewhere
	newDate := date.Format("02/01/2006")
	v, _ := ef.app.gui.View(FORM_VIEW)
	buffer := v.BufferLines()
	cX := len(buffer[len(buffer)-1])
	cY := len(buffer) - 1

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
	ef.entry.Date = newDate
}

func (ef *EntryForm) GetEntry() database.Entry {
	v, _ := ef.app.gui.View(FORM_VIEW)
	ef.app.gui.SetCurrentView(FORM_VIEW)
	v.Clear()

	e := &ef.entry

	date := ef.app.date
	e.Date = date.Format("02/01/2006")
	fmt.Fprintf(v, "Date: %s\n", date.Format(DISPLAY_DATE_FORMAT))

	v.Editable = true

	fmt.Fprintf(v, "Item: ")
	ef.getItem()
	item := ef.app.item
	e.Item = item
	fmt.Fprintln(v, item)

	LeaveItem := true
	if LeaveItem {
		fmt.Fprintf(v, "Logging for date range...\n")
		dateRange :=  ef.GetDateRange()

		//MOVE BELOW TO LOG
		fmt.Fprintf(v, "First day of %s: %s\n", item, dateRange.from.Format(DATE_FORMAT))
		fmt.Fprintf(v, "Last day of %s: %s\n", item, dateRange.to.Format(DATE_FORMAT))
	}
	
	fmt.Fprintf(v, "Hours: ")
	hours := ef.getHours()
	e.Hours = hours

	v.Editable = false
	return *e
}