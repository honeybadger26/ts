package entryform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jroimartin/gocui"
)

const (
	MSG_SUBMIT      = "\x1b[0;32m  Successfully submitted %d hour(s) for %s\x1b[0m"
	SAVED_LOGS_FILE = "data/savedlogs.json"
)

type Form struct {
	gui *gocui.Gui

	itemToSend           chan string
	HandleEntrySubmitted func()
}

func NewForm(g *gocui.Gui, itemToSend chan string) *Form {
	f := &Form{}
	f.gui = g
	f.itemToSend = itemToSend

	maxX, maxY := g.Size()
	v, _ := g.SetView("entryform", 0, 0, maxX/2-1, maxY/2-1)
	v.Title = "New Entry"
	return f
}

type EntryData struct {
	Date  string
	Item  string
	Hours int
}

func (f *Form) getEntryInfo() *EntryData {
	entry := &EntryData{}
	v, _ := f.gui.View("entryform")
	v.Clear()

	t := time.Now()
	entry.Date = t.Format("02/01/2006")
	fmt.Fprintf(v, "Date: %s\n", entry.Date)

	fmt.Fprintf(v, "Item: ")
	itemView := &ItemView{}
	entry.Item = <-itemView.GetItem(f.gui, f.itemToSend)
	fmt.Fprintf(v, "%s\n", entry.Item)

	fmt.Fprintf(v, "Hours: ")
	hoursView := &HoursView{}
	entry.Hours = <-hoursView.GetHours(f.gui)
	fmt.Fprintf(v, "%d\n", entry.Hours)

	return entry
}

func (f *Form) AddEntry() {
	entry := f.getEntryInfo()
	var entries []EntryData

	if _, err := os.Stat(SAVED_LOGS_FILE); os.IsNotExist(err) {
		os.OpenFile(SAVED_LOGS_FILE, os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		data, _ := ioutil.ReadFile(SAVED_LOGS_FILE)
		json.Unmarshal(data, &entries)
	}

	entries = append(entries, *entry)
	file, _ := os.OpenFile(SAVED_LOGS_FILE, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()

	writedata, _ := json.Marshal(entries)
	file.WriteString(string(writedata))

	f.HandleEntrySubmitted()
}
