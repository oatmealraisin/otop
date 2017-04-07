package otop

import gc "github.com/rthornton128/goncurses"

const (
	colorDefault int16 = 0
	colorLow     int16 = 1
	colorMed     int16 = 2
	colorHigh    int16 = 3
	colorSelect  int16 = 4
	colorTab     int16 = 5
	colorHeader  int16 = 6
	colorWarn    int16 = 8
	colorError   int16 = 9
)

// TODO: use a config file to set this
const (
	quitKey       = 'q'
	helpKey       = 'H'
	leftKey       = 'h'
	rightKey      = 'l'
	upKey         = 'k'
	doubleUpKey   = 'K'
	downKey       = 'j'
	doubleDownKey = 'J'
	sortKey       = 's'
	explainKey    = '?'
	editKey       = 'e'
	selectKey     = gc.KEY_RETURN
)

var (
	ResourceModeController map[gc.Key](func(o *Otop) error)
	OverviewModeController map[gc.Key](func(o *Otop) error)

	Tabs = []func(*gc.Window) *Tab{
		NewPodsTab,
		NewBuildsTab,
		NewDeploymentConfigsTab,
		NewRoutesTab,
		NewServicesTab,
	}
)

func init() {
	OverviewModeController = map[gc.Key](func(o *Otop) error){
		quitKey: exitFunction,
		helpKey: showHelp,
		'r':     modeSwitchResource,
	}

	OverviewMode = Mode{
		ActiveTab:  0,
		Controller: OverviewModeController,
		Tabs:       []*Tab{},
	}

	ResourceModeController = map[gc.Key](func(o *Otop) error){
		quitKey:      exitFunction,
		helpKey:      showHelp,
		leftKey:      moveTabLeft,
		rightKey:     moveTabRight,
		'o':          modeSwitchOverview,
		gc.KEY_LEFT:  moveTabLeft,
		gc.KEY_RIGHT: moveTabRight,
	}

	ResourceMode = Mode{
		ActiveTab:  0,
		Controller: ResourceModeController,
		Tabs:       []*Tab{},
	}
}
