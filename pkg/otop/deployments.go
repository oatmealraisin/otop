package otop

import (
	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	gc "github.com/rthornton128/goncurses"
)

func NewDeploymentsTab(w *gc.Window) *Tab {
	e := []map[string]string{}
	panel := gc.NewPanel(w)
	var t *Tab
	maxY, maxX := w.MaxYX()
	subWindow := w.Sub(maxY-1, maxX-1, 2, 0)

	w.Clear()
	w.ColorOn(colorHeader)
	w.HLine(0, 0, ' ', maxX)
	w.MovePrint(0, 0, " Name")
	w.ColorOff(colorHeader)

	t = &Tab{
		Panel:   panel,
		name:    "Deployments",
		entries: e,
		Redraw: func() error {
			subWindow.Clear()
			subMaxY, _ := subWindow.MaxYX()
			for i, entry := range e {
				if i >= subMaxY {
					return nil
				}

				subWindow.MovePrint(i, 0, " "+entry["NAME"])
				i++
			}

			if err := w.Touch(); err != nil {
				return err
			}

			return nil
		},
		Update: func(o *gopenshift.OpenShift) error {
			deployments, err := o.GetDeployments()
			if err != nil {
				return err
			}

			e = []map[string]string{}
			for _, deployment := range deployments {
				e = append(e, map[string]string{
					"NAME": deployment.Name,
				})
			}

			return nil
		},
	}

	return t
}
