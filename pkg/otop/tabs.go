package otop

import (
	"fmt"

	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	gc "github.com/rthornton128/goncurses"
)

// A Tab is a self contained Panel for displaying information. Most panels
// display a specific resource, but this isn't a requirement.
type Tab struct {
	*gc.Panel

	// The name of the tab, printed in the Tab row
	name string
	// The index of the entry that is currently selected.
	activeEntry int
	// entries contains each element of the tab for display
	entries []map[string]string
	// Initialize contains logic for a Tab that is just started. Normally, this
	// just initializes the Header bar with the proper columns, then calls
	// Update and Redraw for the first time.
	Initialize func() error
	// Update is called periodically to refresh information in the tab. This is
	// where a tab should retrieve information to the server, and sort it into
	// entries.
	Update func(*gopenshift.OpenShift) error
	// The redraw function is called every frame, and allows us to keep up
	// with current updates. Normally, this function refreshes a subwindow
	// containing the entries.
	Redraw func() error
}

// Name is a getter for Tab.name
func (t *Tab) Name() string {
	return t.name
}

// Utility function for making sure all of the tab names fit at the top of the
// screen
func formatTitles(titles []string) []string {
	maxSize := 0
	for _, title := range titles {
		if len(title) > maxSize {
			maxSize = len(title)
		}
	}

	for i, title := range titles {
		// NOTE: This formatting looks better, although it makes more sense to
		// NOTE: make the tabs the same size.
		//	titles[i] = fmt.Sprintf(" %-[1]*[2]s", maxSize, title)
		titles[i] = fmt.Sprintf(" %s   ", title)
	}

	return titles
}

// TODO: There's a lot of boilerplate for update functions, let's find a way to
// DRY it
//func genericUpdate(mapper func(*runtime.Object) map[string]string, object *runtime.Object) error {
//	return nil
//}
