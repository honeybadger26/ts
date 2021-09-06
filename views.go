package main

const (
	FORM_VIEW    = "form"
	ITEM_VIEW    = "form.item"
	HOURS_VIEW   = "form.hours"
	HELP_VIEW    = "help"
	ENTRIES_VIEW = "entries"
	INFO_VIEW    = "info"
	LOGGER_VIEW  = "logger"
)

var MAIN_VIEWS = []string{
	FORM_VIEW,
	HELP_VIEW,
	ENTRIES_VIEW,
	INFO_VIEW,
	LOGGER_VIEW,
}

type viewProps struct {
	title          string
	frame          bool
	editable       bool
	x0, y0, x1, y1 float32
}

var VIEW_PROPS = map[string]viewProps{
	FORM_VIEW: {
		title:    "New Entry",
		frame:    true,
		editable: true,
		x0:       0.0,
		y0:       0.0,
		x1:       0.5,
		y1:       0.5,
	},
	ITEM_VIEW: {
		title:    "",
		frame:    true,
		editable: true,
		x0:       -1.0,
		y0:       -1.0,
		x1:       -1.0,
		y1:       -1.0,
	},
	HELP_VIEW: {
		title:    "",
		frame:    false,
		editable: false,
		x0:       0.0,
		y0:       0.5,
		x1:       0.5,
		y1:       1.0,
	},
	ENTRIES_VIEW: {
		title:    "Entries",
		frame:    true,
		editable: false,
		x0:       0.5,
		y0:       0.0,
		x1:       1.0,
		y1:       0.333,
	},
	INFO_VIEW: {
		title:    "Info",
		frame:    true,
		editable: false,
		x0:       0.5,
		y0:       0.333,
		x1:       1.0,
		y1:       0.666,
	},
	LOGGER_VIEW: {
		title:    "Log",
		frame:    true,
		editable: false,
		x0:       0.5,
		y0:       0.666,
		x1:       1.0,
		y1:       1.0,
	},
}