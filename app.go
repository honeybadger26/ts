package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
)

var HELP_TEXT = map[string]string{
	FORM_VIEW: "" +
		"<Ctrl-j> OR <Down> Select next item\n" +
		"<Ctrl-k> OR <Up> Select previous item\n" +
		"<Enter> Confirm selected item\n" +
		"<Alt-l> Next category\n" +
		"<Alt-h> Previous category",
	"APP": "" +
		"<Ctrl-d> Next day\n" +
		"<Ctrl-u> Previous day\n" +
		"<Ctrl-c> Quit",
}

type App struct {
	gui *gocui.Gui
	db  *Database
	ef  *EntryForm

	date time.Time
	item string
}

func NewApp(g *gocui.Gui) *App {
	app := &App{}
	app.gui = g
	app.db = &Database{}

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

	return app
}

func (app *App) setupKeyBindings() {
	app.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})

	app.gui.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		app.changeDate(app.date.AddDate(0, 0, 1))
		return nil
	})

	app.gui.SetKeybinding("", gocui.KeyCtrlU, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		app.changeDate(app.date.AddDate(0, 0, -1))
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

		db := Database{}
		entries := db.getEntries(app.date)
		cols, rows := v.Size()
		padding := cols/3 - 1

		for i, e := range entries {
			rowText := fmt.Sprintf("%-*s %-*s %*d", padding, e.Date, padding, e.Item, padding, e.Hours)

			if i%2 == 0 {
				rowText = fmt.Sprintf("\x1b[0;33m%s\x1b[0m", rowText)
			}

			fmt.Fprintln(v, rowText)
		}

		for len(v.BufferLines()) < rows {
			fmt.Fprintln(v, "")
		}

		pretext := "Total"
		totalHours := db.getTotalHours(app.date)
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
	db := Database{}
	item := db.getItem(app.item)

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
			fmt.Fprintf(v, "Total Hours:    %f\n", item.TotalHours)
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

		if helpText, ok := HELP_TEXT[view]; ok {
			fmt.Fprintln(v, helpText)
		}

		fmt.Fprintf(v, HELP_TEXT["APP"])

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
	if app.db.entryExists(e.Date, e.Item) {
		if e.Hours == 0 {
			msg = fmt.Sprintf("Removing entry: ")
		} else {
			msg = fmt.Sprintf("Updating entry: ")
		}
	} else {
		msg = fmt.Sprintf("Sumbitting new entry: ")
	}

	app.db.saveEntry(e)
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
