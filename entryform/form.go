package entryform

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jroimartin/gocui"
)

const (
	MSG_SUBMIT = "\x1b[0;32mSuccessfully submitted %d hour(s) for %s\x1b[0m"
)

type Form struct {
	gui *gocui.Gui

	HandleItemChange     func(string)
	HandleEntrySubmitted func()
}

func NewForm(g *gocui.Gui) *Form {
	f := &Form{}
	f.gui = g

	maxX, maxY := g.Size()
	v, _ := g.SetView("entryform", 0, 0, maxX/2-1, maxY/2-1)
	v.Title = "New Entry"
	return f
}

func (f *Form) getEntryInfo() (string, int) {
	v, _ := f.gui.View("entryform")

	itemView := &ItemView{}
	itemView.HandleItemChange = f.HandleItemChange
	item := <-itemView.GetItem(f.gui)
	v.Clear()
	fmt.Fprintf(v, "Item: %s\n", item)

	hoursView := &HoursView{}
	hours := <-hoursView.GetHours(f.gui)
	fmt.Fprintf(v, "Hours: %d\n", hours)

	return item, hours
}

func (f *Form) AddEntry() {
	t := time.Now()
	item, hours := f.getEntryInfo()

	line := fmt.Sprintf("%s, %s, %d\n", t.Format("02/01/2006"), item, hours)

	file, err := os.OpenFile("data/savedlogs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()

	if _, err := file.WriteString(line); err != nil {
		log.Panicln(err)
	}

	v, _ := f.gui.View("entryform")
	_, rows := v.Size()
	for i := 0; i < rows-3; i++ {
		fmt.Fprintln(v, "")
	}
	message := fmt.Sprintf(MSG_SUBMIT, hours, item)
	fmt.Fprintln(v, message)

	f.HandleEntrySubmitted()
}
