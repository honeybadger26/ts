package editmode

import (
	"fmt"
	"log"
	"os/exec"
	"time"
	"ts/database"
	"ts/viewmode"
	"ts/whiteboard"

	"github.com/jroimartin/gocui"
)

// should move this to a constants file with views.go
const (
	HELP_TEXT = "" +
		"<F1> Show/hide help\n" +
		"<Up> Select previous item\n" +
		"<Down> Select next item\n" +
		"<Tab> Next category\n" +
		"<Enter> Confirm selected item\n" +
		"<Alt-Left> Previous day\n" +
		"<Alt-Right> Next day\n" +
		"<Ctrl-T> Go to today\n" +
		"<Ctrl-W> Go to weekly view\n" +
		"<Ctrl-X> Quit and sign out of Whiteboard\n" +
		"<Ctrl-C> Quit"
)

type App struct {
	gui *gocui.Gui
	db  *database.Database
	ef  *EntryForm
	va  *viewmode.ViewApp

	date         time.Time
	item         string
	showHelpText bool
}

func NewEditApp(g *gocui.Gui) *App {
	app := &App{}
	app.gui = g
	app.db = &database.Database{}

	app.setupKeyBindings()
	app.setupViews()

	app.date = time.Now()
	app.printEntries()

	app.changeItem("")

	app.printHelp(FORM_VIEW)
	go func() {
		for {
			app.addNewEntry()
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

	app.gui.SetKeybinding("", gocui.KeyCtrlW, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if app.va == nil {
			app.va = viewmode.NewViewApp(g, false)
		} else {
			app.va.Destroy()
			app.va = nil
			g.SetCurrentView(FORM_VIEW)
		}
		return nil
	})

	app.gui.SetKeybinding(FORM_VIEW, gocui.KeyF1, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		app.showHelpText = !app.showHelpText
		viewHelp := VIEW_PROPS[HELP_VIEW]
		viewForm := VIEW_PROPS[FORM_VIEW]
		if app.showHelpText {
			viewHelp.y0 = 0.5
			viewForm.y1 = 0.5
		} else {
			app.log("Help section hidden. Press F1 to unhide.")
			viewHelp.y0 = 1.0
			viewForm.y1 = 1.0
		}
		VIEW_PROPS[HELP_VIEW] = viewHelp
		VIEW_PROPS[FORM_VIEW] = viewForm
		app.setupViews()
		app.ef.updateItemView()
		return nil
	})

	app.gui.SetKeybinding(FORM_VIEW, gocui.KeyCtrlJ, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		item := app.db.GetItem(app.item)
		app.log("Opening " + item.URL + " in browser... ")
		var err = exec.Command("rundll32", "url.dll,FileProtocolHandler", item.URL).Start()
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	app.gui.SetKeybinding(FORM_VIEW, gocui.KeyCtrlX, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		go func() {
			whiteboard.NewWhiteboardHelper(g, FORM_VIEW)
		}()
		return nil
	})
}

func (app *App) setupViews() {
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
	item := app.db.GetItem(app.item)

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
		fmt.Fprintf(v, "Description:    %s\n", item.Description)
		if item.Size != "" {
			fmt.Fprintf(v, "Size:           %s\n", item.Size)
		}
		if item.TotalHours != -1 {
			fmt.Fprintf(v, "Total Hours:    %.2f\n", item.TotalHours+float32(app.db.GetHoursLogged(item.Name)))
		}
		if item.URL != "" {
			fmt.Fprintf(v, "URL:            %s\n", item.URL)
		}

		return nil
	})
}

func (app *App) changeItem(item string) {
	app.item = item
	app.printItemInfo()
}

func (app *App) printHelp(view string) {
	app.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(HELP_VIEW)

		if err != nil {
			return err
		}

		v.Clear()
		fmt.Fprintf(v, HELP_TEXT)

		_, rows := v.Size()
		for len(v.BufferLines()) < rows {
			v.SetCursor(0, 0)
			v.EditNewLine()
		}

		return nil
	})
}

func (app *App) addNewEntry() {
	app.ef = NewEntryForm(app)
	e := app.ef.GetEntry()

	entryStr := fmt.Sprintf("%s - %s - %d hours", e.Date, e.Item, e.Hours)
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
