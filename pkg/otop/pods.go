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
	subWindow := w.Sub(maxY-1, maxX-1, 2, 0)

	// See Tab struct in pkg/otop/tabs.go for more information
	t = &Tab{
		Panel:   panel,
		name:    "Pods",
		entries: e,
		Initialize: func() error {
			w.Clear()
			w.ColorOn(colorHeader)
			w.HLine(0, 0, ' ', maxX)
			w.MovePrint(0, 0, " Name")
			w.MovePrint(0, 15, " Phase")
			w.ColorOff(colorHeader)
			return nil
		},
		Redraw: func() error {
			subWindow.Clear()
			subMaxY, _ := subWindow.MaxYX()
			for i, entry := range e {
				if i >= subMaxY {
					return nil
				}

				subWindow.MovePrint(i, 0, " "+entry["NAME"])
				subWindow.MovePrint(i, 15, " "+entry["PHASE"])
				i++
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
					"NAME":  pod.Name,
					"PHASE": string(pod.Status.Phase),
				})
			}

			return nil
		},
	}

	return t
}
