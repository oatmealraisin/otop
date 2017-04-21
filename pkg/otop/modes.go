package otop

import (
	gc "github.com/rthornton128/goncurses"
)

type Mode struct {
	// The index of the current tab the user is viewing when in this mode
	ActiveTab int

	// The tabs that a user can view in this mode.
	Tabs []*Tab

	// The Controller is a map of keys to functions that handle the keypresses.
	// This cleans up the massive case structure that we used to use for
	// keypresses. The current bindings can be found in defs.go
	Controller map[gc.Key](func(o *Otop) error)
}

var (
	OverviewMode Mode
	ResourceMode Mode
)

func modeSwitchOverview(o *Otop) error {
	o.Mode = &OverviewMode
	o.moveTab(0)
	return nil
}

func modeSwitchResource(o *Otop) error {
	o.Mode = &ResourceMode
	o.moveTab(0)
	return nil
}

func modeSwitchLogs(o *Otop) error {
	return nil
}
