package viewmode

import "ts/viewmanager"

const (
	WEEK_VIEW  = "vm_week"
	HELP_VIEW  = "vm_help"
	INFO_VIEW  = "vm_info"
	BLANK_VIEW = "vm_blank"
)

var MAIN_VIEWS = []string{
	WEEK_VIEW,
	HELP_VIEW,
	INFO_VIEW,
}

var VIEW_PROPS = map[string]viewmanager.ViewProps{
	WEEK_VIEW: {
		Title:      "",
		Frame:      true,
		Editable:   false,
		Wrap:       true,
		Autoscroll: false,
		X0:         0.0,
		Y0:         0.0,
		X1:         1.0,
		Y1:         0.7,
	},
	HELP_VIEW: {
		Title:      "",
		Frame:      false,
		Editable:   false,
		Wrap:       true,
		Autoscroll: false,
		X0:         0.0,
		Y0:         0.7,
		X1:         0.5,
		Y1:         1.0,
	},
	INFO_VIEW: {
		Title:      "",
		Frame:      true,
		Editable:   false,
		Wrap:       true,
		Autoscroll: false,
		X0:         0.5,
		Y0:         0.7,
		X1:         1.0,
		Y1:         1.0,
	},
	BLANK_VIEW: {
		Title:      "",
		Frame:      false,
		Editable:   false,
		Wrap:       true,
		Autoscroll: true,
		X0:         -0.1,
		Y0:         -0.1,
		X1:         1.1,
		Y1:         1.1,
	},
}
