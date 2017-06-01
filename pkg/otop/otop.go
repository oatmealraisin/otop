package otop

import (
	"fmt"
	"io"
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

	// The list of tabs for ResourceMode
	Tabs []*Tab

	// A panel that is preloaded to quickly show what the current keybindings are
	HelpPanel *gc.Panel

	// Our OpenShift client that is used to interact with the API
	OpenShift *gopenshift.OpenShift

	User string

	// The current mode that the user is in. See pkg/otop/modes.go for more info
	Mode *Mode

	// The writer that otop will write logs to, since we can't just write them
	// to STDOUT meaningfully
	w io.Writer
}

func exitFunction(o *Otop) error { return &ExitApplication{} }

// Run initializes otop, then listens for key events. Each key contains a
// function for the resulting action, which may return an error.
func (o *Otop) Run(w io.Writer) error {
	if w != nil {
		o.w = w
		fmt.Fprint(o.w, "Initialized log writer")
	}

	if err := o.init(); err != nil {
		return err
	}

	exitChan := make(chan bool)
	defer func() { exitChan <- true }()

	// Set up a consistent refresh rate
	// TODO: Make a command line argument
	frameRate := int64(10)
	go refresh(o, exitChan, frameRate)

	for {
		// To see all keybindings, check out pkg/otop/defs.go
		keyEvent := gc.Key(o.GetChar())
		execKey := o.Mode.Controller[keyEvent]

		//if execKey == nil {
		//	execKey := o.Tabs[o.activeTab].Events[keyEvent]
		//}

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

// TODO: Handle window resizing
func refresh(o *Otop, exit chan bool, frameRate int64) {
	tick := time.NewTicker(time.Second / time.Duration(frameRate)).C
	for {
		select {
		case <-exit:
			return
		case <-tick:
			if o.w != nil {
				fmt.Fprint(o.w, "TICK")
			}
			o.Mode.Tabs[o.Mode.ActiveTab].Top()
			o.Mode.Tabs[o.Mode.ActiveTab].Update(o.OpenShift)
			o.Mode.Tabs[o.Mode.ActiveTab].Redraw()

			gc.UpdatePanels()
			gc.Update()
		}
	}
}

// Utility method for initialization steps
// NOTE: For now, this method doesn't take very long and doesn't effect start up
// time. However, if it starts to take longer, we may want to look into changing
// this to only initialize the first tab opened.
func (o *Otop) init() error {
	user, err := o.OpenShift.WhoAmI()
	if err != nil {
		return err
	}

	project, err := o.OpenShift.Project()
	if err != nil {
		return err
	}

	if o.Mode == nil {
		o.Mode = &ResourceMode
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
	gc.InitPair(colorHigh, gc.C_RED, -1)
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

	if err := o.printFooter(user, project); err != nil {
		return err
	}

	// TODO: Figure out how to best initialize tabs into modes
	maxY, maxX := o.MaxYX()

	for _, tab := range Tabs {
		if win, err := gc.NewWindow(maxY-footerHeight-displayTop, maxX-displayLeft, displayTop, displayLeft); err == nil {
			t := tab(win)
			// Update all tabs so that the initial tab-switch doesn't take forever
			t.Update(o.OpenShift)
			o.Mode.Tabs = append(o.Mode.Tabs, t)
		}
	}

	ovWin, _ := gc.NewWindow(maxY-footerHeight-displayTop, maxX-displayLeft, displayTop, displayLeft)
	OverviewMode.Tabs = append(OverviewMode.Tabs, NewOverviewTab(ovWin))

	// Draw the initial tab
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
// TODO: Display only visible tabs on smaller screens
func (o *Otop) moveTab(amount int) error {
	// for shorthand
	mode := o.Mode

	// Calculate the new position of the activeTab
	if mode.ActiveTab+amount < 0 {
		amount += mode.ActiveTab
		amount++
		mode.ActiveTab = len(mode.Tabs) - 1
	}

	if mode.ActiveTab+amount > len(mode.Tabs)-1 {
		amount -= (len(mode.Tabs) - mode.ActiveTab - 1)
		amount--
		mode.ActiveTab = 0
	}

	mode.ActiveTab += amount

	return o.printTabs()
}
