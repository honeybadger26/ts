package viewmode

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"

	"ts/database"
)

// repeated code - form.go
const (
	DISPLAY_DATE_FORMAT = "Mon 02 / Jan 01 / 2006"
)

type ViewApp struct {
	gui       *gocui.Gui
	db        *database.Database
	startDate time.Time
	endDate   time.Time
}

func NewViewApp(g *gocui.Gui) *ViewApp {
	va := &ViewApp{}

	va.gui = g
	va.db = &database.Database{}
	va.setDate(time.Now())
	va.setupKeyBindings()

	return va
}

func (va *ViewApp) setupKeyBindings() {
	va.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})

	va.gui.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		va.setDate(va.startDate.AddDate(0, 0, -7))
		return nil
	})

	va.gui.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		va.setDate(va.startDate.AddDate(0, 0, 7))
		return nil
	})

	va.gui.SetKeybinding("", gocui.KeyCtrlT, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		va.setDate(time.Now())
		return nil
	})
}

func (va *ViewApp) setDate(date time.Time) {
	va.startDate = date.Truncate(24 * time.Hour)
	for va.startDate.Weekday() != time.Monday {
		va.startDate = va.startDate.AddDate(0, 0, -1)
	}

	va.endDate = date.Truncate(24 * time.Hour)
	for va.endDate.Weekday() != time.Sunday {
		va.endDate = va.endDate.AddDate(0, 0, 1)
	}

	va.setupViews()
}

func (va *ViewApp) setupViews() {
	maxX, maxY := va.gui.Size()
	x0, y0, x1, y1 := 0, 0, -1, maxY-1
	numDays := 0

	for d := va.startDate; d.Before(va.endDate); d = d.AddDate(0, 0, 1) {
		numDays++
	}

	width := (maxX - 1) / numDays

	for i, d := 0, va.startDate; i <= numDays; i, d = i+1, d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Saturday {
			y1 = ((maxY - 1) / 2) - 1
		}

		if d.Weekday() != time.Sunday {
			x0 = x1 + 1
			x1 = x0 + width - 1
		} else {
			y0 = y1 + 1
			y1 = maxY - 1
		}

		if x1 >= maxX {
			x1 = maxX - 1
		}

		name := fmt.Sprintf("Day%d", i)
		v, err := va.gui.SetView(name, x0, y0, x1, y1)

		if err == nil {
			v.Clear()
		} else if err != gocui.ErrUnknownView {
			log.Panicln(err)
		}

		v.Title = d.Format(DISPLAY_DATE_FORMAT)
		v.Wrap = true
		v.Editable = false
		v.Frame = true

		// repeated code (app.go - printEntries)
		entries := va.db.GetEntries(d)
		cols, rows := v.Size()

		for i, e := range entries {
			hoursStr := fmt.Sprintf("%d", e.Hours)
			rowText := fmt.Sprintf("%-*s%s", cols-len(hoursStr)-1, e.Item, hoursStr)

			if i%2 == 0 {
				rowText = fmt.Sprintf("\x1b[0;33m%s\x1b[0m", rowText)
			}

			fmt.Fprintln(v, rowText)
		}

		// repeated code (app.go - printEntries)
		for len(v.BufferLines()) < rows {
			fmt.Fprintln(v, "")
		}

		pretext := "Total"
		totalHours := va.db.GetTotalHours(d)
		hoursStr := fmt.Sprintf("%d", totalHours)
		fmt.Fprintf(v, "\x1b[0;32m%-*s%s\x1b[0m", cols-len(hoursStr)-1, pretext, hoursStr)
	}
}
