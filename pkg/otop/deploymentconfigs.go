package otop

import (
	"fmt"
	"strconv"

	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	gc "github.com/rthornton128/goncurses"
)

func NewDeploymentConfigsTab(w *gc.Window) *Tab {
	e := []map[string]string{}
	panel := gc.NewPanel(w)
	var t *Tab
	maxY, maxX := w.MaxYX()
	subWindow := w.Sub(maxY-1, maxX-1, 2, 0)
	separators := []int{
		0,  // the start
		58, // Name
		10, // Replicas
	}

	w.Clear()
	w.ColorOn(colorHeader)
	w.HLine(0, 0, ' ', maxX)
	w.MovePrint(0, 0, " Name")
	w.MovePrint(0, separators[1], " Replicas")
	w.ColorOff(colorHeader)

	t = &Tab{
		Panel:   panel,
		name:    "DeploymentConfigs",
		entries: e,
		Redraw: func() error {
			subWindow.Clear()
			subMaxY, _ := subWindow.MaxYX()
			for i, entry := range e {
				if i >= subMaxY {
					return nil
				}
				subWindow.MovePrint(i, 0, fmt.Sprintf(" %s", entry["NAME"]))
				subWindow.MovePrint(i, separators[1], fmt.Sprintf(" %s/%s", entry["REPLICAS_READY"], entry["REPLICAS_DESIRED"]))
				i++
			}

			if err := w.Touch(); err != nil {
				return err
			}

			return nil
		},
		Update: func(o *gopenshift.OpenShift) error {
			dcs, err := o.GetDeploymentConfigs()
			if err != nil {
				return err
			}

			e = []map[string]string{}
			for _, dc := range dcs {
				e = append(e, map[string]string{
					"NAME":             dc.Name,
					"REPLICAS_DESIRED": strconv.Itoa(int(dc.Status.Replicas)),
					"REPLICAS_READY":   strconv.Itoa(int(dc.Status.ReadyReplicas)),
					"PAUSED":           strconv.FormatBool(dc.Spec.Paused),
				})
			}

			return nil
		},
	}

	return t
}
