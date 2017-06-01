package otop

import (
	"fmt"

	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	gc "github.com/rthornton128/goncurses"
	kapi "k8s.io/kubernetes/pkg/api"
)

func NewPodsTab(w *gc.Window) *Tab {
	panel := gc.NewPanel(w)
	var t *Tab
	maxY, maxX := w.MaxYX()
	subWindow := w.Sub(maxY-2, maxX-1, 2, 0)

	pods := []*kapi.Pod{}

	separators := []int{
		0,  // the start
		10, // Status
		21, // Time Alive
	}

	w.Clear()
	w.ColorOn(colorHeader)
	w.HLine(0, 0, ' ', maxX)
	w.MovePrint(0, separators[0], " Status")
	w.MovePrint(0, separators[1], " Uptime")
	w.MovePrint(0, separators[2], " Name")
	w.ColorOff(colorHeader)

	// See Tab struct in pkg/otop/tabs.go for more information
	t = &Tab{
		Panel: panel,
		name:  "Pods",
		Redraw: func() error {

			subWindow.Clear()
			subMaxY, _ := subWindow.MaxYX()
			for i, pod := range pods {

				if i >= subMaxY {
					return nil
				}

				switch string(pod.Status.Phase) {
				case string(kapi.PodSucceeded):
					subWindow.AttrOn(gc.A_DIM)
				case string(kapi.PodFailed):
					subWindow.ColorOn(colorHigh)
					subWindow.AttrOn(gc.A_BOLD)
				}

				subWindow.MovePrint(i, separators[0], fmt.Sprintf(" %s", string(pod.Status.Phase)))
				if pod.Status.StartTime != nil {
					subWindow.MovePrint(i, separators[1], fmt.Sprintf(" %s", formatTimeAlive(pod.Status.StartTime)))
				}
				subWindow.MovePrint(i, separators[2], fmt.Sprintf(" %s", pod.Name))
				subWindow.AttrOff(gc.A_DIM)
				subWindow.AttrOff(gc.A_BOLD)
				subWindow.ColorOff(colorHigh)
			}

			if err := w.Touch(); err != nil {
				return err
			}

			return nil
		},
		// This should give us an array of Pods, but Go has dumb versioning and
		// vendors are dumb
		Update: func(o *gopenshift.OpenShift) error {
			var err error
			pods, err = o.GetPods()
			return err
		},
	}

	return t
}
