package editmode

import (
	"fmt"
	"log"
	"os/exec"
	"time"
	"ts/database"
	"ts/viewmanager"
	"ts/viewmode"

	"github.com/jroimartin/gocui"
)

// should move this to a constants file with views.go
var HELP_TEXT = []string{
	"<F1> Show/hide help",
	" ",
	"<Ctrl-T> Go to today",
	"<Alt-Left> Previous day",
	"<Alt-Right> Next day",
	" ",
	"<Tab> Next category",
	"<Up> Select previous item",
	"<Down> Select next item",
	" ",
	"<Enter> Confirm selected item",
	"<Ctrl-D> Log Date Range (Leave / Public Holiday)",
	"<Ctrl-Z> Cancel latest entry",
	" ",
	"<Ctrl-W> Go to weekly view",
	"<Ctrl-J> Open URL in browser",
	"<Ctrl-X> Quit and sign out of Whiteboard",
	"<Ctrl-C> Quit",
}

const (
	FULL_DAY = 8
	SUNDAY   = 0
	SATURDAY = 6
)

type App struct {
	gui    *gocui.Gui
	db     *database.Database
	psqldb *database.PsqlInterface
	ef     *EntryForm

	date         time.Time
	item         string
	showHelpText bool
}

func NewEditApp(g *gocui.Gui) *App {
	app := &App{}
	app.gui = g
	app.db = &database.Database{}
	app.psqldb = database.NewPsqlInterface()

	app.setupKeyBindings()
	app.setupViews()

	app.date = time.Now()
	app.printEntries()

	app.changeItem("")

	app.printHelp(FORM_VIEW)
	go func() {
		for {
			app.ef = NewEntryForm(app)
			app.addNewEntry(app.ef.GetEntries())
		}
	}()

	app.showHelpText = true

	return app
}

func (app *App) setupKeyBindings() {
	app.gui.SetKeybinding(FORM_VIEW, gocui.KeyArrowLeft, gocui.ModAlt, func(g *gocui.Gui, v *gocui.View) error {
		app.changeDate(app.date.AddDate(0, 0, -1))
		return nil
	})

	app.gui.SetKeybinding(FORM_VIEW, gocui.KeyArrowRight, gocui.ModAlt, func(g *gocui.Gui, v *gocui.View) error {
		app.changeDate(app.date.AddDate(0, 0, 1))
		return nil
	})

	app.gui.SetKeybinding(FORM_VIEW, gocui.KeyCtrlT, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		app.changeDate(time.Now())
		return nil
	})

	app.gui.SetKeybinding(FORM_VIEW, gocui.KeyCtrlW, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		g.Cursor = false

		va := viewmode.NewViewApp(g, app.date, false)
		va.Callback = func() {
			app.changeDate(va.CurrentDate)
			g.SetCurrentView(FORM_VIEW)
			g.Cursor = true
		}

		return nil
	})

	app.gui.SetKeybinding(FORM_VIEW, gocui.KeyF1, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		app.showHelpText = !app.showHelpText
		app.setHelpVisible(app.showHelpText)
		if !app.showHelpText {
			app.log("Help section hidden. Press F1 to unhide.")
		}
		return nil
	})

	app.gui.SetKeybinding(FORM_VIEW, gocui.KeyCtrlJ, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		// item := app.db.GetItem(app.item)
		item := app.psqldb.GetItem(app.item)
		app.log("Opening " + item.URL.String + " in browser... ")
		var err = exec.Command("rundll32", "url.dll,FileProtocolHandler", item.URL.String).Start()
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	app.gui.SetKeybinding(FORM_VIEW, gocui.KeyCtrlX, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		go func() {
			NewWhiteboardHelper(g, FORM_VIEW)
		}()
		return nil
	})

	app.gui.SetKeybinding(FORM_VIEW, gocui.KeyCtrlZ, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if app.db.EntryCount() != 0 {
			var latestEntry database.Entry = app.db.GetLatestEntry()
			latestEntry.Hours = 0
			var entrySlice []database.Entry
			entrySlice = append(entrySlice, latestEntry)
			app.addNewEntry(entrySlice)
		}
		return nil
	})
}

func (app *App) setupViews() {
	for _, n := range MAIN_VIEWS {
		viewmanager.SetupView(app.gui, n, VIEW_PROPS[n])
	}
}

func (app *App) printEntries() {
	app.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(ENTRIES_VIEW)

		if err != nil {
			return err
		}

		v.Clear()

		entries := app.db.GetEntries(app.date)
		cols, rows := v.Size()
		padding := cols/2 - 1

		for i, e := range entries {
			rowText := fmt.Sprintf("%-*s %*d", padding, e.Item, padding, e.Hours)

			if i%2 == 0 {
				rowText = fmt.Sprintf("\x1b[0;33m%s\x1b[0m", rowText)
			}

			fmt.Fprintln(v, rowText)
		}

		for len(v.BufferLines()) < rows {
			fmt.Fprintln(v, "")
		}

		pretext := "Total"
		totalHours := app.db.GetTotalHours(app.date)
		fmt.Fprintf(v, "\x1b[0;32m%s %*d\x1b[0m", pretext, cols-len(pretext)-2, totalHours)

		return nil
	})
}

func (app *App) changeDate(date time.Time) {
	app.date = date
	app.ef.SetDate(date)
	app.printEntries()
}

func (app *App) printItemInfo() {
	// item := app.db.GetItem(app.item)
	item := app.psqldb.GetItem(app.item)

	app.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(INFO_VIEW)

		if err != nil {
			return err
		}

		v.Clear()

		if item == nil {
			return nil
		}

		fmt.Fprintf(v, "Name:           %s\n", item.Name)
		if item.Description.Valid {
			fmt.Fprintf(v, "Description:    %s\n", item.Description.String)
		}
		if item.Size.Valid {
			fmt.Fprintf(v, "Size:           %s\n", item.Size.String)
		}
		if item.URL.Valid {
			fmt.Fprintf(v, "URL:            %s\n", item.URL.String)
		}

		return nil
	})
}

func (app *App) changeItem(item string) {
	app.item = item
	app.printItemInfo()
}

func (app *App) setHelpVisible(visible bool) {
	viewHelp := VIEW_PROPS[HELP_VIEW]
	viewForm := VIEW_PROPS[FORM_VIEW]

	var boundary float32 = 1.0
	if visible {
		_, maxY := app.gui.Size()
		boundary = float32(maxY-len(HELP_TEXT)) / float32(maxY)
	}

	viewHelp.Y0 = boundary
	viewForm.Y1 = boundary
	VIEW_PROPS[HELP_VIEW] = viewHelp
	VIEW_PROPS[FORM_VIEW] = viewForm

	app.setupViews()
	app.ef.updateItemView()
}

func (app *App) printHelp(view string) {
	app.gui.Update(func(g *gocui.Gui) error {
		app.setHelpVisible(true)
		v, err := g.View(HELP_VIEW)

		if err != nil {
			return err
		}

		v.Clear()
		for _, line := range HELP_TEXT {
			fmt.Fprintln(v, line)
		}

		return nil
	})
}

func (app *App) addNewEntry(entrySlice []database.Entry) {
	for _, e := range entrySlice {
		entryStr := fmt.Sprintf("%s - %s - %d hours", e.Date.Format(DATE_FORMAT), e.Item, e.Hours)
		var msg string
		if app.db.EntryExists(e.Date, e.Item) {
			if e.Hours == 0 {
				msg = fmt.Sprintf("Removing entry: ")
			} else {
				msg = fmt.Sprintf("Updating entry: ")
			}
		} else {
			msg = fmt.Sprintf("Sumbitting new entry: ")
		}

		app.db.SaveEntry(e)
		app.log(msg + entryStr)
	}
	app.printEntries()
}

func (app *App) log(text string) {
	app.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(LOGGER_VIEW)

		if err != nil {
			return err
		}

		fmt.Fprintln(v, text)

		_, vHeight := v.Size()
		if len(v.BufferLines()) > vHeight+1 {
			originx, originy := v.Origin()
			v.SetOrigin(originx, originy+1)
		}

		return nil
	})
}
