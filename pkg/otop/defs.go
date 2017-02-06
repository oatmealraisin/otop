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
	Control = map[gc.Key](func(o *Otop) error){
		quitKey:      exitFunction,
		helpKey:      showHelp,
		leftKey:      moveTabLeft,
		rightKey:     moveTabRight,
		gc.KEY_LEFT:  moveTabLeft,
		gc.KEY_RIGHT: moveTabRight,
	}

	Tabs = []func(*gc.Window) *Tab{
		NewPodsTab,
		NewDeploymentsTab,
		NewBuildsTab,
		NewRoutesTab,
		NewServicesTab,
		NewDeploymentConfigsTab,
	}
)
