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

	HELP_TEXT_STANDALONE = "" +
		"<Left> Previous week\n" +
		"<Right> Next week\n" +
		"<Ctrl-t> Go to today\n" +
		"<Ctrl-c> Quit"

	HELP_TEXT_EMBEDDED = "" +
		"<Left> Previous week\n" +
		"<Right> Next week\n" +
		"<Ctrl-t> Go to today\n" +
		"<Ctrl-w> Go to edit mode\n" +
		"<Ctrl-c> Quit"
)

type ViewApp struct {
	gui *gocui.Gui
	db  *database.Database

	standalone bool

	// stores the names of the days that are on the week view
	dayViews  []string
	startDate time.Time
	endDate   time.Time
}

func NewViewApp(g *gocui.Gui, standalone bool) (app *ViewApp) {
	app = &ViewApp{}

	app.gui = g
	app.db = &database.Database{}
	app.standalone = standalone

	app.setupKeyBindings()
	app.setDate(time.Now())
	app.setupViews()

	return
}

func (app *ViewApp) setupKeyBindings() {
	app.gui.SetKeybinding(INFO_VIEW, gocui.KeyArrowLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		app.setDate(app.startDate.AddDate(0, 0, -7))
		app.refreshViews()
		return nil
	})

	app.gui.SetKeybinding(INFO_VIEW, gocui.KeyArrowRight, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		app.setDate(app.startDate.AddDate(0, 0, 7))
		app.refreshViews()
		return nil
	})

	app.gui.SetKeybinding(INFO_VIEW, gocui.KeyCtrlT, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		app.setDate(time.Now())
		app.refreshViews()
		return nil
	})
}

func (app *ViewApp) setDate(date time.Time) {
	app.startDate = date.Truncate(24 * time.Hour)
	for app.startDate.Weekday() != time.Monday {
		app.startDate = app.startDate.AddDate(0, 0, -1)
	}

	app.endDate = date.Truncate(24 * time.Hour)
	for app.endDate.Weekday() != time.Sunday {
		app.endDate = app.endDate.AddDate(0, 0, 1)
	}
}

func (app *ViewApp) printHelp() {
	app.gui.Update(func(g *gocui.Gui) error {
		// repeated code - app.go, printHelp
		v, err := g.View(HELP_VIEW)

		if err != nil {
			return err
		}

		v.Clear()

		if app.standalone {
			fmt.Fprintf(v, HELP_TEXT_STANDALONE)
		} else {
			fmt.Fprintf(v, HELP_TEXT_EMBEDDED)
		}

		// hacky way to get text to be at bottom of view
		_, rows := v.Size()
		for len(v.BufferLines()) < rows {
			v.SetCursor(0, 0)
			v.EditNewLine()
		}

		return nil
	})
}

func (app *ViewApp) setupViews() {
	// repeated code (app.go in editmode)
	maxX, maxY := app.gui.Size()

	for _, name := range MAIN_VIEWS {
		p := VIEW_PROPS[name]
		x0 := int(p.x0 * float32(maxX))
		y0 := int(p.y0 * float32(maxY))
		x1 := int(p.x1*float32(maxX)) - 1
		y1 := int(p.y1*float32(maxY)) - 1

		if !p.frame {
			y0 = y0 - 1
			y1 = y1 + 1
		}

		if v, err := app.gui.SetView(name, x0, y0, x1, y1); err != nil {
			if err != gocui.ErrUnknownView {
				log.Panicln(err)
			}
			v.Title = p.title
			v.Wrap = true
			v.Editable = p.editable
			v.Frame = p.frame
		}
	}

	app.refreshViews()
	app.printHelp()
	app.gui.SetCurrentView(INFO_VIEW)
}

// make this relative to WEEK_VIEW
func (app *ViewApp) refreshViews() {
	g := app.gui
	p := VIEW_PROPS[WEEK_VIEW]
	maxX, maxY := g.Size()
	parentx0 := int(p.x0 * float32(maxX))
	parenty0 := int(p.y0 * float32(maxY))
	parentx1 := int(p.x1*float32(maxX)) - 1
	parenty1 := int(p.y1*float32(maxY)) - 1

	var x0, y0, x1, y1 int
	numDays := 0

	for d := app.startDate; d.Before(app.endDate); d = d.AddDate(0, 0, 1) {
		numDays++
	}

	width := (parentx1 - parentx0 - 1) / numDays

	for i, d := 0, app.startDate; i <= numDays; i, d = i+1, d.AddDate(0, 0, 1) {
		// surely there's a better way to do this?
		if d.Weekday() == time.Monday {
			x0 = parentx0 + 1
			y0 = parenty0 + 1
			x1 = x0 + width - 1
			y1 = parenty1 - 1
		} else if d.Weekday() == time.Saturday {
			x0 = x1 + 1
			x1 = parentx1 - 1
			y1 = parenty0 + 1 + ((parenty1 - parenty0) / 2)
		} else if d.Weekday() == time.Sunday {
			y0 = y1 + 1
			y1 = parenty1 - 1
		} else {
			x0 = x1 + 1
			x1 = x0 + width - 1
		}

		name := fmt.Sprintf("Day%d", i)
		v, err := app.gui.SetView(name, x0, y0, x1, y1)

		if err == nil {
			v.Clear()
		} else if err != gocui.ErrUnknownView {
			log.Panicln(err)
		} else {
			app.dayViews = append(app.dayViews, name)
		}

		v.Title = d.Format(DISPLAY_DATE_FORMAT)
		v.Wrap = true
		v.Editable = false
		v.Frame = true

		// repeated code (app.go - printEntries)
		entries := app.db.GetEntries(d)
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
		totalHours := app.db.GetTotalHours(d)
		hoursStr := fmt.Sprintf("%d", totalHours)
		fmt.Fprintf(v, "\x1b[0;32m%-*s%s\x1b[0m", cols-len(hoursStr)-1, pretext, hoursStr)
	}

	// also refresh info view
}

func (app *ViewApp) Destroy() {
	for _, name := range MAIN_VIEWS {
		app.gui.DeleteView(name)
	}

	for _, name := range app.dayViews {
		app.gui.DeleteView(name)
	}
}
