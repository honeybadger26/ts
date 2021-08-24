package entryform

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	MSG_SUBMIT = "\x1b[0;32mSuccessfully submitted %d hour(s) for %s\x1b[0m"
)

type Form struct {
	gui              *gocui.Gui
	handleItemChange func(string)
}

func (f *Form) Init(g *gocui.Gui, handleItemChange func(string)) {
	f.gui = g
	f.handleItemChange = handleItemChange

	maxX, maxY := g.Size()
	v, _ := g.SetView("entryform", 0, 0, maxX/2-1, maxY/2-1)
	v.Title = "New Entry"
}

func (f *Form) AddEntry() (string, int) {
	v, _ := f.gui.View("entryform")

	itemView := &ItemView{}
	itemView.HandleItemChange = f.handleItemChange
	item := <-itemView.GetItem(f.gui)
	v.Clear()
	fmt.Fprintf(v, "Item: %s\n", item)

	hoursView := &HoursView{}
	hours := <-hoursView.GetHours(f.gui)
	fmt.Fprintf(v, "Hours: %d\n", hours)

	_, rows := v.Size()
	for i := 0; i < rows-3; i++ {
		fmt.Fprintln(v, "")
	}
	message := fmt.Sprintf(MSG_SUBMIT, hours, item)
	fmt.Fprintln(v, message)

	return item, hours
}
