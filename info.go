package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jroimartin/gocui"
)

type ItemData struct {
	Name        string
	Description string
	Size        string
	TotalHours  float32
}

type InfoComponent struct {
	gui *gocui.Gui
}

func NewInfo(g *gocui.Gui) *InfoComponent {
	i := &InfoComponent{}
	i.gui = g

	maxX, maxY := g.Size()

	if v, err := g.SetView("info", maxX/2, maxY/2, maxX-1, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			return nil
		}
		v.Wrap = true
		v.Editable = false
		v.Frame = true
		v.Title = "Item Info"
	}

	return i
}

func (i *InfoComponent) getItemInfo(item string) *ItemData {
	file, err := os.Open("data/items.json")

	if err != nil {
		return nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data ItemData

	decoder.Token()

	for decoder.More() {
		decoder.Decode(&data)
		if data.Name == item {
			break
		}
	}

	return &data
}

func (i *InfoComponent) UpdateInfo(item string) {
	v, err := i.gui.View("info")

	if err != nil {
		return
	}

	v.Clear()

	if item == "" {
		return
	}

	info := i.getItemInfo(item)
	fmt.Fprintf(v, "Name: %s\n", info.Name)
	fmt.Fprintf(v, "Description: %s\n", info.Description)
	fmt.Fprintf(v, "Size: %s\n", info.Size)
	fmt.Fprintf(v, "Total Hours: %f\n", info.TotalHours)
}
