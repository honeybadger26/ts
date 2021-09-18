package editmode

import "ts/viewmanager"

const (
	FORM_VIEW    = "em_form"
	ITEM_VIEW    = "em_form.item"
	HOURS_VIEW   = "em_form.hours"
	HELP_VIEW    = "em_help"
	ENTRIES_VIEW = "em_entries"
	INFO_VIEW    = "em_info"
	LOGGER_VIEW  = "em_logger"
)

var MAIN_VIEWS = []string{
	FORM_VIEW,
	HELP_VIEW,
	ENTRIES_VIEW,
	INFO_VIEW,
	LOGGER_VIEW,
}

var VIEW_PROPS = map[string]viewmanager.ViewProps{
	FORM_VIEW: {
		Title:      "New Entry",
		Frame:      true,
		Editable:   true,
		Wrap:       true,
		Autoscroll: false,
		X0:         0.0,
		Y0:         0.0,
		X1:         0.5,
		Y1:         0.5,
	},
	ITEM_VIEW: {
		Title:      "",
		Frame:      true,
		Editable:   true,
		Wrap:       true,
		Autoscroll: false,
		X0:         -1.0,
		Y0:         -1.0,
		X1:         -1.0,
		Y1:         -1.0,
	},
	HELP_VIEW: {
		Title:      "",
		Frame:      false,
		Editable:   false,
		Wrap:       true,
		Autoscroll: false,
		X0:         0.0,
		Y0:         0.4,
		X1:         0.5,
		Y1:         1.0,
	},
	ENTRIES_VIEW: {
		Title:      "Entries",
		Frame:      true,
		Editable:   false,
		Wrap:       true,
		Autoscroll: false,
		X0:         0.5,
		Y0:         0.333,
		X1:         1.0,
		Y1:         0.80,
	},
	INFO_VIEW: {
		Title:      "Info",
		Frame:      true,
		Editable:   false,
		Wrap:       true,
		Autoscroll: false,
		X0:         0.5,
		Y0:         0.0,
		X1:         1.0,
		Y1:         0.333,
	},
	LOGGER_VIEW: {
		Title:      "Log",
		Frame:      true,
		Editable:   false,
		Wrap:       true,
		Autoscroll: true,
		X0:         0.5,
		Y0:         0.80,
		X1:         1.0,
		Y1:         1.0,
	},
}
