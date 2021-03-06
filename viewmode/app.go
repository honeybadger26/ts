package viewmode

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"

	"ts/database"
	"ts/viewmanager"
)

// repeated code - form.go
const (
	DISPLAY_DATE_FORMAT = "Mon 02 / Jan 01 / 2006"
	WEEK_HOUR_LIMIT     = 40
	DAY_HOUR_LIMIT      = 24
	ITEM_LOWER_LIMIT    = 7
)

var HELP_TEXT = []string{
	"<F1> Show/hide help",
	"",
	"<Left> Previous day",
	"<Right> Next day",
	"<Ctrl-T> Go to today",
	"",
	"<Alt-Left> Previous week",
	"<Alt-Right> Next week",
	"<Ctrl-H> Show/hide weekend",
	"",
	"<Ctrl-C> Quit",
}

type ViewApp struct {
	gui *gocui.Gui
	db  *database.Database

	showHelp    bool
	showWeekend bool
	standalone  bool

	// stores the names of the days that are on the week view
	dayViews []string

	CurrentDate time.Time
	Callback    func()
}

func NewViewApp(g *gocui.Gui, date time.Time, standalone bool) (app *ViewApp) {
	app = &ViewApp{}

	app.gui = g
	app.db = &database.Database{}

	app.showHelp = true
	app.showWeekend = true
	app.standalone = standalone
	app.CurrentDate = date

	app.setupKeyBindings()
	app.setupViews()
	app.printHelp()
	g.SetCurrentView(INFO_VIEW)

	return
}

type DateChange int

const (
	dcDay DateChange = iota
	dcWeek
)

func (app *ViewApp) handleDateChange(forward bool, changeType DateChange) {
	var numDays int

	if changeType == dcDay {
		numDays = 1
		if !app.showWeekend && ((forward && app.CurrentDate.Weekday() == time.Friday) ||
			(!forward && app.CurrentDate.Weekday() == time.Monday)) {
			numDays = 3
		}
	} else {
		numDays = 7
	}

	if !forward {
		numDays *= -1
	}

	app.CurrentDate = app.CurrentDate.AddDate(0, 0, numDays)
	app.refreshViews()
}

func (app *ViewApp) setupKeyBindings() {
	if !app.standalone {
		app.gui.SetKeybinding(INFO_VIEW, gocui.KeyCtrlW, gocui.ModNone, app.destroy)
	}

	app.gui.SetKeybinding(INFO_VIEW, gocui.KeyF1, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		app.showHelp = !app.showHelp
		app.printHelp()
		return nil
	})

	app.gui.SetKeybinding(INFO_VIEW, gocui.KeyArrowLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		app.handleDateChange(false, dcDay)
		return nil
	})

	app.gui.SetKeybinding(INFO_VIEW, gocui.KeyArrowRight, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		app.handleDateChange(true, dcDay)
		return nil
	})

	app.gui.SetKeybinding(INFO_VIEW, gocui.KeyArrowLeft, gocui.ModAlt, func(g *gocui.Gui, v *gocui.View) error {
		app.handleDateChange(false, dcWeek)
		return nil
	})

	app.gui.SetKeybinding(INFO_VIEW, gocui.KeyArrowRight, gocui.ModAlt, func(g *gocui.Gui, v *gocui.View) error {
		app.handleDateChange(true, dcWeek)
		return nil
	})

	app.gui.SetKeybinding(INFO_VIEW, gocui.KeyCtrlT, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		app.CurrentDate = time.Now()
		if time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday {
			app.showWeekend = true
		}
		app.refreshViews()
		return nil
	})

	app.gui.SetKeybinding(INFO_VIEW, gocui.KeyCtrlH, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		app.showWeekend = !app.showWeekend
		app.refreshViews()
		return nil
	})
}

func (app *ViewApp) getDateRange() (start time.Time, end time.Time) {
	// start = app.CurrentDate.Truncate(24 * time.Hour)
	start = app.CurrentDate
	for start.Weekday() != time.Monday {
		start = start.AddDate(0, 0, -1)
	}

	// end = app.CurrentDate.Truncate(24 * time.Hour)
	end = app.CurrentDate
	for end.Weekday() != time.Sunday {
		end = end.AddDate(0, 0, 1)
	}

	return
}

func (app *ViewApp) printHelp() {
	// repeated code app.go in editmode setHelpVisible
	viewWeek := VIEW_PROPS[WEEK_VIEW]
	viewHelp := VIEW_PROPS[HELP_VIEW]
	viewInfo := VIEW_PROPS[INFO_VIEW]

	helpLines := HELP_TEXT
	if !app.standalone {
		helpLines = append(HELP_TEXT[:len(HELP_TEXT)-1], HELP_TEXT[len(HELP_TEXT)-2:]...)
		helpLines[len(helpLines)-2] = "<Ctrl-W> Go to edit view"
	}

	if app.showHelp {
		_, maxY := app.gui.Size()
		y := float32(maxY-len(helpLines)) / float32(maxY)
		viewWeek.Y1 = y
		viewHelp.Y0 = y
		viewInfo.Y0 = y
		viewInfo.X0 = 0.5
	} else {
		viewInfo.X0 = 0.0
	}

	VIEW_PROPS[WEEK_VIEW] = viewWeek
	VIEW_PROPS[HELP_VIEW] = viewHelp
	VIEW_PROPS[INFO_VIEW] = viewInfo

	app.setupViews()
	app.refreshViews()

	if !app.showHelp {
		return
	}

	app.gui.Update(func(g *gocui.Gui) error {
		// repeated code - app.go, printHelp
		v, err := g.View(HELP_VIEW)

		if err != nil {
			return err
		}

		v.Clear()

		for _, line := range helpLines {
			fmt.Fprintln(v, line)
		}

		return nil
	})
}

func (app *ViewApp) setupViews() {
	g := app.gui

	if !app.standalone {
		viewmanager.SetupView(g, BLANK_VIEW, VIEW_PROPS[BLANK_VIEW])
	}

	for _, n := range MAIN_VIEWS {
		viewmanager.SetupView(g, n, VIEW_PROPS[n])
	}
}

// make this relative to WEEK_VIEW
func (app *ViewApp) refreshViews() {
	g := app.gui
	p := VIEW_PROPS[WEEK_VIEW]

	parentx0, parenty0, parentx1, parenty1 :=
		viewmanager.GetDimensions(g, p.X0, p.Y0, p.X1, p.Y1)

	startDate, endDate := app.getDateRange()

	// width is 6 instead of 7 because weekend days are stacked
	numDays := 7
	width := (parentx1 - parentx0 - 1) / (numDays - 1)
	if !app.showWeekend {
		numDays = 5
		width = (parentx1 - parentx0 - 1) / numDays
	}

	for _, name := range app.dayViews {
		g.DeleteView(name)
	}
	app.dayViews = []string{}

	var x0, y0, x1, y1 int

	for i, d := 0, startDate; i < numDays; i, d = i+1, d.AddDate(0, 0, 1) {
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
		cd := app.CurrentDate

		var v *gocui.View
		var err error
		if d.Year() == cd.Year() && d.YearDay() == cd.YearDay() {
			// doesn't really work for weekends when they are hidden
			v, err = app.gui.SetView(name, x0, y0-1, x1, y1+1)
		} else {
			v, err = app.gui.SetView(name, x0, y0, x1, y1)
		}

		if err == nil {
			v.Clear()
		} else if err != gocui.ErrUnknownView {
			log.Panicln(numDays)
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

		totalColor := "32"
		if totalHours < ITEM_LOWER_LIMIT || totalHours > DAY_HOUR_LIMIT {
			totalColor = "31"
		}
		fmt.Fprintf(v, "\x1b[0;"+totalColor+"m%-*s%s\x1b[0m", cols-len(hoursStr)-1, pretext, hoursStr)
	}

	// also refresh info view
	v, err := g.View(INFO_VIEW)
	if err != nil && err != gocui.ErrUnknownView {
		log.Panicln(err)
	}

	v.Clear()

	totalWeekHours := app.db.GetTotalHoursForRange(startDate, endDate)
	infoText := fmt.Sprintf("Total hours this week: %d\n", totalWeekHours)
	if totalWeekHours < WEEK_HOUR_LIMIT {
		infoText += fmt.Sprintf("\x1b[0;33mHours logged this week is under %d\x1b[0m", WEEK_HOUR_LIMIT)
	}

	fmt.Fprint(v, infoText)
}

func (app *ViewApp) destroy(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView(BLANK_VIEW)

	for _, name := range MAIN_VIEWS {
		g.DeleteView(name)
	}

	for _, name := range app.dayViews {
		g.DeleteView(name)
	}

	app.Callback()
	return nil
}
