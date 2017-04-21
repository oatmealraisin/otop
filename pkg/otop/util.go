package otop

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"k8s.io/kubernetes/pkg/api/unversioned"

	gc "github.com/rthornton128/goncurses"
)

// printTabs
// Utility method for printing tabs. Will throw an error if there are too many
// characters for the window
// tab.
func (o *Otop) printTabs() error {
	maxY, _ := o.MaxYX()
	o.Move(tabTop, 0)

	// Don't bother if the screen is too small
	if maxY <= 5 {
		return nil
	}

	if len(o.Mode.Tabs) == 0 {
		return fmt.Errorf("No tabs")
	}

	if err := o.ClearToEOL(); err != nil {
		return err
	}

	// Print each tab name, highlight the active tab
	for i := 0; i < len(o.Mode.Tabs); i++ {
		if i == o.Mode.ActiveTab {
			o.ColorOn(colorSelect)
			o.Print(fmt.Sprintf(" %s  ", o.Mode.Tabs[i].Name()))
			o.ColorOff(colorSelect)
		} else {
			o.Print(fmt.Sprintf(" %s  ", o.Mode.Tabs[i].Name()))
		}
	}

	o.Mode.Tabs[o.Mode.ActiveTab].Top()

	if err := o.Mode.Tabs[o.Mode.ActiveTab].Redraw(); err != nil {
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

	o.ColorOn(colorTab)
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
	o.ColorOff(colorTab)

	return nil
}

//. TODO: Take into account too many tabs for screen
func moveTabLeft(o *Otop) error {
	return o.moveTab(-1)
}

//. TODO: Take into account too many tabs for screen
func moveTabRight(o *Otop) error {
	return o.moveTab(1)
}

func showHelp(o *Otop) error {
	maxY, maxX := o.MaxYX()

	if maxY > 10 && maxY > 10 {
		w := o.Sub(maxY-10, maxX-10, 5, 5)

		if err := w.Touch(); err != nil {
			return err
		}

		if err := w.Keypad(true); err != nil {
			return err
		}

		w.Border(gc.ACS_VLINE, gc.ACS_VLINE, gc.ACS_HLINE, gc.ACS_HLINE,
			gc.ACS_ULCORNER, gc.ACS_URCORNER, gc.ACS_LLCORNER, gc.ACS_LRCORNER)
		w.AttrOn(gc.A_UNDERLINE)
		w.MovePrintf(1, 1, "Help ")
		w.AttrOff(gc.A_UNDERLINE)

		maxY, maxX := w.MaxYX()

		for i, s := range []string{
			fmt.Sprintf("%c/%c, Left/Right : Change tab", leftKey, rightKey),
			fmt.Sprintf("%c/%c, Up/Down    : Select resource", upKey, downKey),
			fmt.Sprintf("%c               : Exit otop", quitKey),
			fmt.Sprintf("%c               : Explain resource", explainKey),
		} {
			if i+3 < maxY && len(s)+2 < maxX {
				w.MovePrint(i+3, 2, s)
			}
		}
		w.GetChar()
		w.Clear()
		if err := o.Touch(); err != nil {
			return err
		}

		// TODO: Need to redraw covered things
		o.Tabs[o.Mode.ActiveTab].Update(o.OpenShift)
		o.Tabs[o.Mode.ActiveTab].Redraw()
		gc.UpdatePanels()
		gc.Update()

		if err := w.Delete(); err != nil {
			return err
		}

	}

	return nil
}

var secondsTrimRegexp = regexp.MustCompile(`\.[0-9]+`)

func formatTimeAlive(startTime *unversioned.Time) string {
	//          Uptime
	// 1y300d20h15m10s

	// Unfortunately, kube has it's own time struct..
	parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", startTime.String())
	if err != nil {
		return err.Error()
	}

	return secondsTrimRegexp.ReplaceAllString(time.Since(parsedTime).String(), "")
}

func formatProgressBar(percent float64, length int) string {
	return fmt.Sprintf("[%-*s]", 10, strings.Repeat("|", int(float64(length)*percent)))
}
