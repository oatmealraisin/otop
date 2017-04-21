package otop

import (
	"fmt"
	"strings"

	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	gc "github.com/rthornton128/goncurses"
	kapi "k8s.io/kubernetes/pkg/api"
)

func NewServicesTab(w *gc.Window) *Tab {
	panel := gc.NewPanel(w)
	var t *Tab
	maxY, maxX := w.MaxYX()
	subWindow := w.Sub(maxY-1, maxX-1, 2, 0)
	separators := []int{
		0,  // the start
		13, // Type
		28, // Cluster IP
		40, // Ports
	}

	services := []*kapi.Service{}

	w.Clear()
	w.ColorOn(colorHeader)
	w.HLine(0, 0, ' ', maxX)
	w.MovePrint(0, separators[0], " Type")
	w.MovePrint(0, separators[1], " Cluster IP")
	w.MovePrint(0, separators[2], " Ports")
	w.MovePrint(0, separators[len(separators)-1], " Name")
	w.ColorOff(colorHeader)

	t = &Tab{
		Panel: panel,
		name:  "Services",
		Redraw: func() error {
			subWindow.Clear()
			subMaxY, _ := subWindow.MaxYX()
			for i, service := range services {
				if i >= subMaxY {
					return nil
				}

				subWindow.MovePrint(i, separators[0], fmt.Sprintf(" %s", service.Spec.Type))
				subWindow.MovePrint(i, separators[1], fmt.Sprintf(" %s", service.Spec.ClusterIP))
				subWindow.MovePrint(i, separators[2], fmt.Sprintf(" %s", concatPorts(service.Spec.Ports)))
				subWindow.MovePrint(i, separators[len(separators)-1], fmt.Sprintf(" %s", service.Name))
			}

			if err := w.Touch(); err != nil {
				return err
			}

			return nil
		},
		Update: func(o *gopenshift.OpenShift) error {
			var err error
			services, err = o.GetServices()
			return err
		},
	}

	return t
}

func concatPorts(ports []kapi.ServicePort) string {
	result := ""
	for _, port := range ports {
		result = fmt.Sprintf("%s%d/%s,", result, port.Port, port.Protocol)
	}
	return strings.TrimRight(result, ",")
}
