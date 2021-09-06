package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jroimartin/gocui"
)

func filterItems(g *gocui.Gui, allItems []ItemData, query string) []ItemData {
	filteredItems := []ItemData{}

	for _, item := range allItems {
		regexstr := `(?i)` + query
		match, err := regexp.MatchString(regexstr, item.Name)
		if err != nil {
			// handle error
		}
		if match {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

func updateItemView(g *gocui.Gui, items []ItemData, index int) {
	g.DeleteView(ITEM_VIEW)
	printHelp(g, "")

	if len(items) == 0 {
		return
	}

	printHelp(g, ITEM_VIEW)

	fv, _ := g.View(FORM_VIEW)
	p := VIEW_PROPS[FORM_VIEW]
	maxX, maxY := g.Size()
	x0 := int(p.x0*float32(maxX)) + 1
	y0 := int(p.y0*float32(maxY)) + 1 + len(fv.BufferLines())
	x1 := int(p.x1*float32(maxX)) - 1 - 1
	y1 := int(p.y1*float32(maxY)) - 1 - 1

	if iv, err := g.SetView(ITEM_VIEW, x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return
		}

		iv.Wrap = true
		iv.Editable = VIEW_PROPS[ITEM_VIEW].editable
		iv.Frame = VIEW_PROPS[ITEM_VIEW].frame

		for i, item := range items {
			if i == index {
				fmt.Fprintf(iv, "\x1b[0;34m> %s\x1b[0m\n", item.Name)
			} else {
				fmt.Fprintln(iv, item.Name)
			}
		}
	}
}

func getItem(g *gocui.Gui) string {
	items := []ItemData{}
	file, err := os.Open("data/items.json")

	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data ItemData

	decoder.Token()

	for decoder.More() {
		decoder.Decode(&data)
		items = append(items, data)
	}

	v, _ := g.View(FORM_VIEW)
	buffer := v.BufferLines()
	cX := len(buffer[len(buffer)-1])
	cY := len(buffer) - 1
	v.SetCursor(cX, cY)

	item := make(chan string)

	// put in function
	selectedIndex := -1
	filteredItems := filterItems(g, items, "")
	if len(filteredItems) != 0 {
		selectedIndex = 0
		printItemInfo(g, filteredItems[selectedIndex].Name)
	}
	updateItemView(g, filteredItems, selectedIndex)

	v.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		buf := v.BufferLines()
		line := buf[len(buf)-1]
		oldInput := strings.TrimSpace(line[cX:])

		switch {
		case key == gocui.KeyArrowDown || key == gocui.KeyCtrlJ:
			if selectedIndex < len(filteredItems)-1 {
				selectedIndex++
			}
			updateItemView(g, filteredItems, selectedIndex)
			printItemInfo(g, filteredItems[selectedIndex].Name)
			return
		case key == gocui.KeyArrowUp || key == gocui.KeyCtrlK:
			if selectedIndex > 0 {
				selectedIndex--
			}
			updateItemView(g, filteredItems, selectedIndex)
			printItemInfo(g, filteredItems[selectedIndex].Name)
			return
		case key == gocui.KeyEnter:
			if selectedIndex != -1 {
				g.DeleteView(ITEM_VIEW)
				printHelp(g, "")
				v.SetCursor(cX, cY)
				for range filteredItems[selectedIndex].Name {
					v.EditDelete(false)
				}
				item <- filteredItems[selectedIndex].Name
			}
			return
		case key == gocui.KeyBackspace || key == gocui.KeyBackspace2 || key == gocui.KeyArrowLeft:
			if newCursorX, newCursorY := v.Cursor(); newCursorX == cX && newCursorY == cY {
				return
			}
		}
		gocui.DefaultEditor.Edit(v, key, ch, mod)

		buf = v.BufferLines()
		line = buf[len(buf)-1]
		newInput := strings.TrimSpace(line[cX:])
		if newInput != oldInput {
			selectedIndex = -1
			filteredItems = filterItems(g, items, newInput)
			if len(filteredItems) != 0 {
				selectedIndex = 0
				printItemInfo(g, filteredItems[selectedIndex].Name)
			}
			updateItemView(g, filteredItems, selectedIndex)
		}
	})

	return <-item
}

func getHours(g *gocui.Gui) int {
	v, _ := g.View(FORM_VIEW)
	buffer := v.BufferLines()
	cX := len(buffer[len(buffer)-1])
	cY := len(buffer) - 1
	v.SetCursor(cX, cY)

	hours := make(chan int)

	v.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		switch {
		case ch != 0 && mod == 0 && !unicode.IsNumber(ch):
			return
		case key == gocui.KeyEnter:
			buf := v.BufferLines()
			line := buf[len(buf)-1]
			hoursStr := strings.TrimSpace(line[cX:])
			hoursInt, _ := strconv.Atoi(hoursStr)
			v.SetCursor(cX, cY)
			for range hoursStr {
				v.EditDelete(false)
			}
			hours <- hoursInt
			return
		case key == gocui.KeyBackspace || key == gocui.KeyBackspace2 || key == gocui.KeyArrowLeft:
			if newCursorX, newCursorY := v.Cursor(); newCursorX == cX && newCursorY == cY {
				return
			}
		}
		gocui.DefaultEditor.Edit(v, key, ch, mod)
	})

	return <-hours
}

func addNewEntry(g *gocui.Gui) {
	entry := &EntryData{}

	v, _ := g.View(FORM_VIEW)
	g.SetCurrentView(FORM_VIEW)
	v.Clear()

	t := time.Now()
	entry.Date = t.Format("02/01/2006")
	fmt.Fprintf(v, "Date: %s\n", entry.Date)

	v.Editable = true

	fmt.Fprintf(v, "Item: ")
	entry.Item = getItem(g)
	fmt.Fprintln(v, entry.Item)

	fmt.Fprintf(v, "Hours: ")
	entry.Hours = getHours(g)
	fmt.Fprintln(v, entry.Hours)

	v.Editable = false

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

	msg := fmt.Sprintf("Successfully submitted entry for: %s for %d hours on %s", entry.Item, entry.Hours, entry.Date)
	logger(g, msg)
}
