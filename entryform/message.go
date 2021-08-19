package entryform

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	MSG_SUBMIT = "\x1b[0;32mSuccessfully submitted %d hour(s) for %s"
)

type MessageView struct {
	gui *gocui.Gui
}

func (mv *MessageView) Init(g *gocui.Gui) {
	mv.gui = g
}

func (mv *MessageView) Clear() {
	mv.gui.DeleteView("message")
}

func (mv *MessageView) Show(item string, hours int) {
	maxX, _ := mv.gui.Size()
	message := fmt.Sprintf(MSG_SUBMIT, hours, item)

	v, _ := mv.gui.SetView("message", 0, 0, maxX/2-1, 3)
	v.Frame = false
	fmt.Fprintf(v, message)
}
