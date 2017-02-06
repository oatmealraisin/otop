package otop

import (
	"fmt"

	gc "github.com/rthornton128/goncurses"
)

// printTabs
// Utility method for printing tabs. Will throw an error if there are too many
// characters for the window
// tab.
//
// Errors:
//    WindowTooSmallError
func (o *Otop) printTabs(tabRange []int) error {
	maxY, _ := o.MaxYX()
	o.Move(tabTop, 0)

	if maxY <= 5 {
		return nil
	}

	if len(o.Tabs) == 0 {
		return fmt.Errorf("No tabs")
	}

	for _, tab := range tabRange {
		if tab == o.activeTab {
			o.ColorOn(colorSelect)
			o.Print(fmt.Sprintf(" %s  ", o.Tabs[tab].Name()))
			o.ColorOff(colorSelect)
		} else {
			o.Print(fmt.Sprintf(" %s  ", o.Tabs[tab].Name()))
		}
	}

	o.Tabs[o.activeTab].Top()

	if err := o.Tabs[o.activeTab].Redraw(); err != nil {
		return err
	}

	gc.UpdatePanels()
	gc.Update()

	return nil
}

// printFooter
// Utility method for printing the footer.
//
// Errors:
func (o *Otop) printFooter(user, project string) error {
	maxY, maxX := o.MaxYX()
	if maxY <= 8 {
		return nil
	}

	o.HLine(maxY-footerHeight, 0, ' ', maxX)
	if user == "" && project != "" {
		o.MovePrintf(maxY-footerHeight, 1, "%s", project)
	}

	if user != "" && project == "" {
		o.MovePrintf(maxY-footerHeight, 1, "%s", user)
	}

	if user != "" && project != "" {
		o.MovePrintf(maxY-footerHeight, 1, "%s/%s", user, project)
	}

	helpString := "Press 'H' for controls"
	if (len(helpString) + len(user) + len(project) + 3) < maxX {
		o.MovePrintf(maxY-footerHeight, maxX-1-len(helpString), "%s", helpString)
	}

	return nil
}
