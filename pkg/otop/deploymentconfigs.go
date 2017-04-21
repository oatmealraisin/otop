package otop

import (
	"fmt"

	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	deployapi "github.com/openshift/origin/pkg/deploy/api"
	gc "github.com/rthornton128/goncurses"
)

func NewDeploymentConfigsTab(w *gc.Window) *Tab {
	panel := gc.NewPanel(w)
	var t *Tab
	maxY, maxX := w.MaxYX()
	subWindow := w.Sub(maxY-1, maxX-1, 2, 0)
	separators := []int{
		0,  // the start
		13, // Replicas
	}

	deploymentConfigs := []*deployapi.DeploymentConfig{}

	w.Clear()
	w.ColorOn(colorHeader)
	w.HLine(0, 0, ' ', maxX)
	w.MovePrint(0, separators[0], " Replicas")
	w.MovePrint(0, separators[len(separators)-1], " Name")
	w.ColorOff(colorHeader)

	t = &Tab{
		Panel: panel,
		name:  "DeploymentConfigs",
		Redraw: func() error {
			subWindow.Clear()
			subMaxY, _ := subWindow.MaxYX()
			for i, dc := range deploymentConfigs {
				if i >= subMaxY {
					return nil
				}
				subWindow.MovePrint(i, separators[0], fmt.Sprintf(" %s", formatProgressBar(float64(dc.Status.ReadyReplicas)/float64(dc.Status.Replicas), 10)))
				subWindow.MovePrint(i, separators[len(separators)-1], fmt.Sprintf(" %s", dc.Name))
			}

			if err := w.Touch(); err != nil {
				return err
			}

			return nil
		},
		Update: func(o *gopenshift.OpenShift) error {
			var err error
			deploymentConfigs, err = o.GetDeploymentConfigs()
			return err
		},
	}

	return t
}
