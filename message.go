package main

import "fmt"

const (
	MSG_SUBMIT = "\x1b[0;32mSuccessfully submitted %d hour(s) for %s"
)

type MessageView struct {
	controller *Controller
}

func (mv *MessageView) Init(c *Controller) error {
	mv.controller = c
	return nil
}

func (mv *MessageView) ShowSubmittedMessage() {
	item := mv.controller.itemView.Item
	hours := mv.controller.hoursView.Hours

	g := mv.controller.gui
	maxX, _ := g.Size()

	message := fmt.Sprintf(MSG_SUBMIT, hours, item)

	v, _ := g.SetView("message", 0, 6, maxX/2-1, 9)
	v.Frame = false
	fmt.Fprintf(v, message)
}
