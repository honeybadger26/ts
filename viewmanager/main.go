package viewmanager

import (
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

type ViewProps struct {
	Title          string
	Frame          bool
	Editable       bool
	Wrap           bool
	Autoscroll     bool
	X0, Y0, X1, Y1 float32
}

func GetDimensions(g *gocui.Gui, x0, y0, x1, y1 float32) (absx0, absy0, absx1, absy1 int) {
	maxX, maxY := g.Size()
	absx0 = int(x0 * float32(maxX))
	absy0 = int(y0 * float32(maxY))
	absx1 = int(x1*float32(maxX)) - 1
	absy1 = int(y1*float32(maxY)) - 1
	return
}

func SetupView(g *gocui.Gui, name string, props ViewProps) {
	x0, y0, x1, y1 := GetDimensions(g, props.X0, props.Y0, props.X1, props.Y1)

	if !props.Frame {
		y0 = y0 - 1
		y1 = y1 + 1
	}

	if v, err := g.SetView(name, x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			log.Panicln(err)
		}

		v.Title = props.Title
		v.Frame = props.Frame
		v.Editable = props.Editable
		v.Wrap = props.Wrap
		v.Autoscroll = props.Autoscroll
	}
}

func GetEndPos(v *gocui.View) (x, y int) {
	buffer := v.BufferLines()
	x = len(buffer[len(buffer)-1])
	y = len(buffer) - 1
	return
}

func GetEndString(v *gocui.View, index int) (result string) {
	buf := v.BufferLines()
	line := buf[len(buf)-1]
	result = strings.TrimSpace(line[index:])
	return
}
