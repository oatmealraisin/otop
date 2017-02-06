package otop

import (
	"fmt"
	"time"

	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	gc "github.com/rthornton128/goncurses"
)

const (
	headerTop    = 0
	headerLeft   = 0
	tabTop       = 0
	tabLeft      = 0
	displayTop   = 1
	displayLeft  = 0
	footerHeight = 1
)

type ExitApplication struct{}

func (e *ExitApplication) Error() string { return "" }

// The main struct for the program.  Otop is a set of settings and variables
// that every method uses.  It is a window, so any window method can be used
// on it.  This means that a controller can start multiple instances of Otop
// in one terminal.
type Otop struct {
	*gc.Window

	Tabs []*Tab

	OpenShift *gopenshift.OpenShift
	HelpPanel *gc.Panel

	// The Controller is a map of keys to functions that handle the keypresses.
	// This cleans up the massive case structure that we used to use for
	// keypresses. The current bindings can be found in defs.go
	Controller map[gc.Key](func(o *Otop) error)

	User string

	// The index of the currently shown tab
	activeTab    int
	resourceMode bool
}

func exitFunction(o *Otop) error { return &ExitApplication{} }

// Run initializes otop, then listens for key events. Each key contains a
// function for the resulting action, which may return an error.
func (o *Otop) Run() error {

	if err := o.init(); err != nil {
		return err
	}

	exitChan := make(chan bool)
	defer func() { exitChan <- true }()

	go o.refresh(exitChan)

	for {
		execKey := o.Controller[gc.Key(o.GetChar())]
		if execKey == nil {
			continue
		}

		err := execKey(o)

		if err != nil {
			if _, ok := err.(*ExitApplication); ok {
				return nil
			}
			return err
		}
	}

	return nil
}

func (o *Otop) refresh(exit chan bool) {
	tick := time.NewTicker(time.Second / 10).C
	for {
		select {
		case <-exit:
			return
		case <-tick:
			o.Tabs[o.activeTab].Update(o.OpenShift)
			o.Tabs[o.activeTab].Redraw()
			gc.UpdatePanels()
			gc.Update()
		}
	}
}

// Utility method for initialization steps
func (o *Otop) init() error {
	user, err := o.OpenShift.WhoAmI()
	if err != nil {
		return err
	}

	project, err := o.OpenShift.Project()
	if err != nil {
		return err
	}

	if o.Controller == nil {
		o.Controller = Control
	}

	// Initializing our colors
	// TODO: Handle colorless terminal
	if err := gc.StartColor(); err != nil {
		return err
	}

	// Sets the FG and BG colors to -1
	if err := gc.UseDefaultColors(); err != nil {
		return err
	}

	// When the user presses keys, don't put them in the screen
	gc.Echo(false)
	// No cursor either
	gc.Cursor(0)

	// Initialize colors
	gc.InitPair(colorDefault, -1, -1)
	gc.InitPair(colorLow, gc.C_BLACK, gc.C_WHITE)
	gc.InitPair(colorMed, -1, gc.C_YELLOW)
	gc.InitPair(colorHigh, gc.C_BLACK, gc.C_RED)
	gc.InitPair(colorSelect, -1, gc.C_BLUE)
	gc.InitPair(colorHeader, gc.C_BLACK, gc.C_GREEN)
	gc.InitPair(colorTab, -1, gc.C_BLUE)
	gc.InitPair(colorError, gc.C_BLACK, gc.C_WHITE)

	o.SetBackground(gc.ColorPair(colorDefault))
	o.Keypad(true)

	//helpWindow := o.Sub(10, 20, 4, 8)
	//hp := gc.NewPanel(helpWindow)
	//o.HelpPanel = hp
	//o.HelpPanel.Hide()

	o.ColorOn(colorTab)
	if err := o.printFooter(user, project); err != nil {
		return err
	}
	o.ColorOff(colorTab)

	maxY, maxX := o.MaxYX()
	for _, tab := range Tabs {
		if win, err := gc.NewWindow(maxY-footerHeight-displayTop, maxX-displayLeft, displayTop, displayLeft); err == nil {
			t := tab(win)
			t.Update(o.OpenShift)
			t.Initialize()
			o.Tabs = append(o.Tabs, t)
		}
	}

	// Draw all tabs
	if err := o.moveTab(0); err != nil {
		return err
	}

	return nil
}

// moveTab moves the tab a certain amount in a direction. A negative number will
// move left, positive will move right. On hitting the end of the visible tabs,
// it will either
//   a) Hide the furthest tab and show a closer tab
//   b) Move to the first/last tab
func (o *Otop) moveTab(amount int) error {
	// Calculate the new position of the activeTab
	if o.activeTab+amount < 0 {
		amount += o.activeTab
		amount++
		o.activeTab = len(o.Tabs) - 1
	}

	if o.activeTab+amount > len(o.Tabs)-1 {
		amount -= (len(o.Tabs) - o.activeTab - 1)
		amount--
		o.activeTab = 0
	}

	o.activeTab += amount

	// TODO: Display only visible tabs on smaller screens
	selection := []int{}
	for i, _ := range o.Tabs {
		selection = append(selection, i)
	}

	return o.printTabs(selection)
	// Get the set of displayed tabs
	//maxY, maxX := o.MaxYX()

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
		o.Tabs[o.activeTab].Update(o.OpenShift)
		o.Tabs[o.activeTab].Redraw()
		gc.UpdatePanels()
		gc.Update()

		if err := w.Delete(); err != nil {
			return err
		}

	}

	return nil
}
