package otop

import "fmt"

// printTabs
// Utility method for printing tabs. Will throw an error if there are too many
// characters for the window
// tab.
//
// Errors:
//    WindowTooSmallError
func (o *Otop) printTabs(tabRange []int) error {
	maxY, maxX := o.MaxYX()
	o.Move(tabTop, 1)

	if maxY <= 5 {
		return nil
	}

	if len(o.Tabs) == 0 {
		return fmt.Errorf("No tabs")
	}

	if tabRange[0] > 0 && maxX > len(o.Tabs[tabRange[0]].Name())+10 {
		o.Print("... ")
	}

	for _, tab := range tabRange {
		_, x := o.YX()
		if x+len(o.Tabs[tab].Name()) > maxX-1 {
			return fmt.Errorf("Too many tabs for window size.")
		}
		if tab == o.activeTab {
			o.ColorOn(colorSelect)
			o.Print(o.Tabs[tab].Name())
			o.ColorOff(colorSelect)
		} else {
			o.Print(o.Tabs[tab].Name())
		}
	}

	if tabRange[len(tabRange)-1] < len(o.Tabs)-1 && maxX > len(o.Tabs[tabRange[0]].Name())+10 {
		o.Print(" ...")
	}

	o.MovePrint(10, 10, o.OpenShift.GetPods())

	o.Refresh()
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
