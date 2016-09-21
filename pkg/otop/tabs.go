package otop

import (
	gos "github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	gc "github.com/rthornton128/goncurses"
)

type Tab struct {
	*gc.Panel

	name        string
	activeEntry int
	entries     Entries
	Resource    string

	Update func(o *gos.OpenShift) error
}

func NewTab(n string, w *gc.Window) *Tab {
	return &Tab{
		Panel: gc.NewPanel(w),
		name:  n,
		Update: func(o *gos.OpenShift) error {
			return nil
		},
	}
}
func (t *Tab) Name() string {
	return t.name
}

type Entry map[string]string

// TODO: Stub
func (e *Entry) String() string {
	return ""
}

type Entries []Entry

// TODO: Stub
func (e *Entries) String() string {
	return ""
}
