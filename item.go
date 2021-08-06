package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jroimartin/gocui"
)

// TODO: Rename to ItemView
type ItemView struct {
	controller *Controller

	// List of all items user can log hours for
	items []string

	// List of items filtered out after searching
	filteredItems []string

	// Currently selected item. Shows '>' next to it
	selectedItem int

	// Item that hours will be logged for. Chosen by user when pressing 'Enter' on
	// selected item
	Item string
}

func (fv *ItemView) Init(c *Controller) error {
	fv.controller = c
	fv.selectedItem = -1

	if err := fv.importItems(); err != nil {
		return err
	}

	if err := fv.setupView(); err != nil {
		return err
	}

	return nil
}

func (fv *ItemView) importItems() error {
	file, err := os.Open("data/items")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fv.items = append(fv.items, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (fv *ItemView) handleEnter() {
	if fv.selectedItem == -1 {
		return
	}

	fv.Item = fv.filteredItems[fv.selectedItem]

	g := fv.controller.gui
	g.DeleteView("item.results")

	maxX, _ := g.Size()
	v, _ := g.SetView("item", 0, 0, maxX/2-1, 2)
	v.Editable = false
	v.Clear()
	fmt.Fprintln(v, "\x1b[0;34m"+fv.Item)

	fv.controller.hoursView.Focus()
}

func (fv *ItemView) searchItems(query string) {
	fv.selectedItem = -1
	fv.filteredItems = nil

	for _, item := range fv.items {
		regexstr := `(?i)` + query
		match, err := regexp.MatchString(regexstr, item)
		if err != nil {
			// handle error
		}
		if match {
			fv.filteredItems = append(fv.filteredItems, item)
		}
	}

	if len(fv.filteredItems) != 0 {
		fv.selectedItem = 0
	}
}

func (fv *ItemView) updateView() {
	g := fv.controller.gui
	maxX, _ := g.Size()
	numResults := len(fv.filteredItems)

	if numResults == 0 {
		g.DeleteView("item.results")
		g.SetView("item", 0, 1, maxX/2-1, 3)
		return
	}

	g.SetView("item", 0, 1, maxX/2-1, 5+numResults)
	v, _ := g.SetView("item.results", 1, 3, maxX/2-2, 4+numResults)

	v.Clear()
	for i, item := range fv.filteredItems {
		if i == fv.selectedItem {
			fmt.Fprintln(v, `> `+item)
		} else {
			fmt.Fprintln(v, item)
		}
	}
}

func (fv *ItemView) setupView() error {
	g := fv.controller.gui
	maxX, _ := g.Size()

	if v, err := g.SetView("item", 0, 1, maxX/2-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		v.Editable = true
		v.Frame = true
		v.Title = "Item name"

		if _, err := g.SetCurrentView("item"); err != nil {
			return err
		}

		v.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			queryBefore := strings.TrimSpace(v.Buffer())
			selectedItemBefore := fv.selectedItem

			switch {
			case ch != 0 && mod == 0:
				v.EditWrite(ch)
			case key == gocui.KeySpace:
				v.EditWrite(' ')
			case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
				v.EditDelete(true)
			case key == gocui.KeyArrowDown:
				if fv.selectedItem != len(fv.filteredItems)-1 {
					fv.selectedItem++
				}
			case key == gocui.KeyArrowUp:
				if fv.selectedItem != 0 {
					fv.selectedItem--
				}
			case key == gocui.KeyEnter:
				fv.handleEnter()
				return
			}

			if query := strings.TrimSpace(v.Buffer()); query != queryBefore {
				fv.searchItems(query)
				fv.updateView()
			}

			if fv.selectedItem != selectedItemBefore {
				fv.updateView()
			}
		})
	}

	return nil
}

func (fv *ItemView) VaildItemSelected() bool {
	return fv.selectedItem != -1
}
