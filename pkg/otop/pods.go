package otop

import (
	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	gc "github.com/rthornton128/goncurses"
)

func NewPodsTab(w *gc.Window) *Tab {
	e := []map[string]string{}
	panel := gc.NewPanel(w)
	var t *Tab
	maxY, maxX := w.MaxYX()
	subWindow := w.Sub(maxY-2, maxX-1, 2, 0)

	separators := []int{
		0,  // the start
		40, // Name
		60, // Status
	}

	w.Clear()
	w.ColorOn(colorHeader)
	w.HLine(0, 0, ' ', maxX)
	w.MovePrint(0, separators[0], " Name")
	w.MovePrint(0, separators[1], " Phase")
	w.ColorOff(colorHeader)

	// See Tab struct in pkg/otop/tabs.go for more information
	t = &Tab{
		Panel:   panel,
		name:    "Pods",
		entries: e,
		Redraw: func() error {

			subWindow.Clear()
			subMaxY, _ := subWindow.MaxYX()
			for i, entry := range e {

				if i >= subMaxY {
					return nil
				}

				subWindow.MovePrint(i, separators[0], " "+entry["NAME"])
				subWindow.MovePrint(i, separators[1], " "+entry["PHASE"])
				i++
			}

			if err := w.Touch(); err != nil {
				return err
			}

			return nil
		},
		Update: func(o *gopenshift.OpenShift) error {
			pods, err := o.GetPods()
			if err != nil {
				return err
			}

			e = []map[string]string{}
			for _, pod := range pods {
				e = append(e, map[string]string{
					"NAME":       pod.Name,
					"PHASE":      string(pod.Status.Phase),
					"START_TIME": pod.Status.StartTime.String(),
				})
			}

			return nil
		},
	}

	return t
}
