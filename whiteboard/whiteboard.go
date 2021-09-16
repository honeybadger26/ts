package whiteboard

import (
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
)

type WhiteboardHelper struct {
	gui          *gocui.Gui
	width        int
	previousView string
}

const (
	PADDING_LINES = 3
	VIEW_NAME     = "wb_main"
)

var INITIAL_TEXT = []string{
	"Are you sure you want to quit and sign out of Whiteboard?",
	"<Ctrl-x> Yes    <Ctrl-z> No",
}
var PROGRESS_TEXT = []string{
	"Signing you out of Whiteboard...",
}
var DONE_TEXT = []string{
	"Signing you out of Whiteboard...done",
	"Quitting...",
}

func NewWhiteboardHelper(g *gocui.Gui, previousView string) (wh *WhiteboardHelper) {
	wh = &WhiteboardHelper{}

	wh.gui = g
	wh.previousView = previousView

	wh.setupKeyBindings()
	wh.setupView(INITIAL_TEXT)

	return
}

func (wh *WhiteboardHelper) getViewContent(text []string) (content []string) {
	for i := 0; i < PADDING_LINES; i++ {
		content = append(content, " ")
	}

	for _, t := range text {
		leftPadding := 0
		rightPadding := wh.width - len(t)
		paddingStr := ""
		for leftPadding < rightPadding {
			leftPadding++
			rightPadding--
			paddingStr += " "
		}
		content = append(content, paddingStr+t)
	}

	for i := 0; i < PADDING_LINES; i++ {
		content = append(content, " ")
	}

	return
}

func (wh *WhiteboardHelper) setupKeyBindings() {
	wh.gui.SetKeybinding(VIEW_NAME, gocui.KeyCtrlX, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		wh.setupView(PROGRESS_TEXT)

		go func() {
			time.Sleep(2 * time.Second)
			wh.setupView(DONE_TEXT)
			time.Sleep(2 * time.Second)
			g.Update(func(g *gocui.Gui) error {
				return gocui.ErrQuit
			})
		}()

		return nil
	})

	wh.gui.SetKeybinding(VIEW_NAME, gocui.KeyCtrlZ, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		g.Cursor = true
		g.DeleteView(VIEW_NAME)
		g.SetCurrentView(wh.previousView)
		return nil
	})
}

func (wh *WhiteboardHelper) setupView(text []string) {
	wh.gui.Update(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()
		height := len(text) + (2 * PADDING_LINES)
		wh.width = maxX / 3

		x0 := int(0.5*float32(maxX)) - (wh.width / 2)
		y0 := int(0.5*float32(maxY)) - (height / 2)
		x1 := int(0.5*float32(maxX)) + (wh.width / 2)
		y1 := y0 + height + 1

		v, err := g.SetView(VIEW_NAME, x0, y0, x1, y1)

		if err == nil {
			v.Clear()
		} else if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		v.Editable = false
		v.Frame = true

		g.Cursor = false
		g.SetViewOnTop(VIEW_NAME)
		g.SetCurrentView(VIEW_NAME)

		for _, t := range wh.getViewContent(text) {
			fmt.Fprintln(v, t)
		}
		return nil
	})
}
