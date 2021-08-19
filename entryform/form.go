package entryform

import "github.com/jroimartin/gocui"

// TODO: Fix this so that don't ned to use controller
// Do the same thing with other componenets if possible?
type Form struct {
	gui *gocui.Gui

	message *MessageView
}

func (f *Form) Init(g *gocui.Gui) {
	f.gui = g
	f.message = &MessageView{}
	f.message.Init(g)
}

func (f *Form) AddEntry() (string, int) {
	f.message.Clear()

	itemView := &ItemView{}
	item := <-itemView.GetItem(f.gui)

	hoursView := &HoursView{}
	hours := <-hoursView.GetHours(f.gui)

	f.message.Show(item, hours)
	return item, hours
}
