package otop

import (
	"fmt"

	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	buildapi "github.com/openshift/origin/pkg/build/api"
	deployapi "github.com/openshift/origin/pkg/deploy/api"
	routeapi "github.com/openshift/origin/pkg/route/api"
	gc "github.com/rthornton128/goncurses"
	kapi "k8s.io/kubernetes/pkg/api"
)

type Application struct {
	Pods             []*kapi.Pod
	DeploymentConfig *deployapi.DeploymentConfig
	Builds           []*buildapi.Build
	Routes           []*routeapi.Route
}

func NewOverviewTab(w *gc.Window) *Tab {
	panel := gc.NewPanel(w)
	var t *Tab
	maxY, maxX := w.MaxYX()
	subWindow := w.Sub(maxY-2, maxX-1, 2, 0)

	applications := map[string]Application{}

	separators := []int{
		0, // the start
	}

	w.Clear()
	w.ColorOn(colorHeader)
	w.HLine(0, 0, ' ', maxX)
	w.MovePrint(0, separators[0], " Application")
	w.ColorOff(colorHeader)

	t = &Tab{
		Panel: panel,
		name:  "Overview",
		Redraw: func() error {
			subWindow.Clear()
			subMaxY, _ := subWindow.MaxYX()
			i := 0
			for name, _ := range applications {

				if i >= subMaxY {
					return nil
				}

				subWindow.MovePrint(i, separators[0], fmt.Sprintf(" %s", name))
				i++
			}

			if err := w.Touch(); err != nil {
				return err
			}

			return nil
		},
		Update: func(o *gopenshift.OpenShift) error {
			applications = make(map[string]Application)

			deploymentConfigs, err := o.GetDeploymentConfigs()
			if err != nil {
				return err
			}

			for _, dc := range deploymentConfigs {
				if label := dc.Labels["app"]; label != "" {
					applications[label] = Application{}
				}
			}

			return nil
		},
	}

	return t
}
