package entryform

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jroimartin/gocui"
)

type ItemView struct {
	gui *gocui.Gui

	// List of all items user can log hours for
	items []string

	// List of items filtered out after searching
	filteredItems []string

	// Currently selected item. Shows '>' next to it
	selectedItem int

	// Item that hours will be logged for. Chosen by user when pressing 'Enter' on
	// selected item
	item chan string
}

func (iv *ItemView) importItems() {
	iv.items = []string{}
	file, err := os.Open("data/items")

	if err != nil {
		// TODO fix this
		// iv.controller.logger.Log("Could not open items file")
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		iv.items = append(iv.items, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		// iv.controller.logger.Log("Could not read items file")
	}
}

func (iv *ItemView) handleEnter() {
	if iv.selectedItem == -1 {
		return
	}

	iv.gui.DeleteView("item.results")
	iv.gui.DeleteView("item")
	iv.item <- iv.filteredItems[iv.selectedItem]
}

func (iv *ItemView) searchItems(query string) {
	iv.filteredItems = nil

	for _, item := range iv.items {
		regexstr := `(?i)` + query
		match, err := regexp.MatchString(regexstr, item)
		if err != nil {
			// handle error
		}
		if match {
			iv.filteredItems = append(iv.filteredItems, item)
		}
	}

	numResults := len(iv.filteredItems)
	if numResults == 0 {
		iv.selectedItem = -1
	} else if iv.selectedItem < 0 || iv.selectedItem >= numResults {
		iv.selectedItem = 0
	}
}

func (iv *ItemView) updateView() {
	g := iv.gui
	maxX, _ := g.Size()

	if iv.selectedItem == -1 {
		g.DeleteView("item.results")
		g.SetView("item", 1, 1, maxX/2-2, 4)
		return
	}

	numResults := len(iv.filteredItems)

	g.SetView("item", 1, 1, maxX/2-2, 5+numResults)
	v, _ := g.SetView("item.results", 2, 3, maxX/2-3, 4+numResults)

	v.Clear()
	for i, item := range iv.filteredItems {
		if i == iv.selectedItem {
			fmt.Fprintln(v, `> `+item)
		} else {
			fmt.Fprintln(v, item)
		}
	}
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
	case key == gocui.KeyArrowDown:
		if iv.selectedItem != len(iv.filteredItems)-1 {
			iv.selectedItem++
		}
	case key == gocui.KeyArrowUp:
		if iv.selectedItem != 0 {
			iv.selectedItem--
		}
	}

	query := strings.TrimSpace(v.Buffer())
	iv.searchItems(query)
	iv.updateView()
}

func (iv *ItemView) GetItem(g *gocui.Gui) chan string {
	iv.gui = g
	iv.selectedItem = -1
	iv.item = make(chan string)

	maxX, _ := g.Size()

	if v, err := g.SetView("item", 1, 1, maxX/2-2, 4); err != nil {
		if err != gocui.ErrUnknownView {
			return nil
		}

		v.Wrap = true
		v.Editable = true
		v.Frame = true
		v.Title = "Item"

		if _, err := g.SetCurrentView("item"); err != nil {
			return nil
		}

		v.Editor = gocui.EditorFunc(iv.editorFunc)

		iv.importItems()
		iv.searchItems("")
		iv.updateView()
	}

	return iv.item
}
