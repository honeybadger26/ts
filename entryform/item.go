package entryform

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jroimartin/gocui"
)

const (
	itemx0 = 1
	itemy0 = 3
)

type ItemData struct {
	Name string
}

type ItemView struct {
	gui *gocui.Gui

	// List of all items user can log hours for
	items []ItemData

	// List of items filtered out after searching
	filteredItems []ItemData

	// Currently selected item. Shows '>' next to it
	selectedItem int

	// Item that hours will be logged for. Chosen by user when pressing 'Enter' on
	// selected item
	item chan string

	// Function that updates the info view
	HandleItemChange func(string)
}

func (iv *ItemView) importItems() {
	file, err := os.Open("data/items.json")

	if err != nil {
		// TODO fix this
		// iv.controller.logger.Log("Could not open items file")
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data ItemData

	decoder.Token()

	for decoder.More() {
		decoder.Decode(&data)
		iv.items = append(iv.items, data)
	}
}

func (iv *ItemView) handleEnter() {
	if iv.selectedItem == -1 {
		return
	}

	iv.gui.DeleteView("item.results")
	iv.gui.DeleteView("item")
	iv.item <- iv.filteredItems[iv.selectedItem].Name
}

func (iv *ItemView) searchItems(query string) {
	iv.filteredItems = nil

	for _, item := range iv.items {
		regexstr := `(?i)` + query
		match, err := regexp.MatchString(regexstr, item.Name)
		if err != nil {
			// handle error
		}
		if match {
			iv.filteredItems = append(iv.filteredItems, item)
		}
	}

	numResults := len(iv.filteredItems)
	g := iv.gui
	maxX, _ := g.Size()
	itemx1 := maxX/2 - 2

	if numResults == 0 {
		iv.selectedItem = -1
		g.DeleteView("item.results")
		g.SetView("item", itemx0, itemy0, itemx1, itemy0+2)
		iv.HandleItemChange("")
		return
	} else if iv.selectedItem < 0 || iv.selectedItem >= numResults {
		iv.selectedItem = 0
	}

	g.SetView("item", itemx0, itemy0, itemx1, itemy0+3+numResults)
	v, _ := g.SetView("item.results", itemx0+1, itemy0+1, itemx1-1, itemy0+3+numResults)

	v.Frame = false

	v.Clear()
	cols, _ := v.Size()
	fmt.Fprintln(v, strings.Repeat("-", cols))

	for i, item := range iv.filteredItems {
		if i == iv.selectedItem {
			fmt.Fprintf(v, "\x1b[0;34m> %s\x1b[0m\n", item.Name)
		} else {
			fmt.Fprintln(v, item.Name)
		}
	}
	iv.HandleItemChange(iv.filteredItems[iv.selectedItem].Name)
}

func (iv *ItemView) editorFunc(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if key == gocui.KeyEnter {
		iv.handleEnter()
		return
	}

	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyArrowDown || key == gocui.KeyCtrlJ:
		if iv.selectedItem != len(iv.filteredItems)-1 {
			iv.selectedItem++
		}
	case key == gocui.KeyArrowUp || key == gocui.KeyCtrlK:
		if iv.selectedItem != 0 {
			iv.selectedItem--
		}
	}

	query := strings.TrimSpace(v.Buffer())
	iv.searchItems(query)
}

func (iv *ItemView) GetItem(g *gocui.Gui) chan string {
	iv.gui = g
	iv.selectedItem = -1
	iv.item = make(chan string)

	maxX, _ := g.Size()

	if v, err := g.SetView("item", itemx0, itemy0, maxX/2-2, itemy0+2); err != nil {
		if err != gocui.ErrUnknownView {
			return nil
		}

		v.Wrap = true
		v.Editable = true
		v.Frame = true

		if _, err := g.SetCurrentView("item"); err != nil {
			return nil
		}

		v.Editor = gocui.EditorFunc(iv.editorFunc)

		iv.importItems()
		iv.searchItems("")
	}

	return iv.item
}
