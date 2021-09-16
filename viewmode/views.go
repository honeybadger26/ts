package viewmode

const (
	WEEK_VIEW = "vm_week"
	HELP_VIEW = "vm_help"
	INFO_VIEW = "vm_info"
)

var MAIN_VIEWS = []string{
	WEEK_VIEW,
	HELP_VIEW,
	INFO_VIEW,
}

type viewProps struct {
	title          string
	frame          bool
	editable       bool
	x0, y0, x1, y1 float32
}

var VIEW_PROPS = map[string]viewProps{
	WEEK_VIEW: {
		title:    "Week",
		frame:    true,
		editable: false,
		x0:       0.0,
		y0:       0.0,
		x1:       1.0,
		y1:       0.7,
	},
	HELP_VIEW: {
		title:    "",
		frame:    false,
		editable: false,
		x0:       0.0,
		y0:       0.7,
		x1:       0.5,
		y1:       1.0,
	},
	INFO_VIEW: {
		title:    "Info",
		frame:    true,
		editable: false,
		x0:       0.5,
		y0:       0.7,
		x1:       1.0,
		y1:       1.0,
	},
}
