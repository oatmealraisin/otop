package otop

import (
	"fmt"

	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	routeapi "github.com/openshift/origin/pkg/route/api"
	gc "github.com/rthornton128/goncurses"
)

func NewRoutesTab(w *gc.Window) *Tab {
	panel := gc.NewPanel(w)
	var t *Tab
	maxY, maxX := w.MaxYX()
	subWindow := w.Sub(maxY-1, maxX-1, 2, 0)

	separators := []int{
		0,
	}

	w.Clear()
	w.ColorOn(colorHeader)
	w.HLine(0, 0, ' ', maxX)
	w.MovePrint(0, 0, " Name")
	w.ColorOff(colorHeader)

	routes := []*routeapi.Route{}

	t = &Tab{
		Panel: panel,
		name:  "Routes",
		Redraw: func() error {
			subWindow.Clear()
			subMaxY, _ := subWindow.MaxYX()
			for i, route := range routes[:subMaxY] {
				if i >= subMaxY {
					return nil
				}

				subWindow.MovePrint(i, separators[len(separators)-1], fmt.Sprintf(" %s", route.Name))
			}

			if err := w.Touch(); err != nil {
				return err
			}

			return nil
		},
		Update: func(o *gopenshift.OpenShift) error {
			var err error
			routes, err = o.GetRoutes()
			return err
		},
	}

	return t
}
