package otop

import (
	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	gc "github.com/rthornton128/goncurses"
)

func NewServicesTab(w *gc.Window) *Tab {
	e := []map[string]string{}
	panel := gc.NewPanel(w)
	var t *Tab
	maxY, maxX := w.MaxYX()
	subWindow := w.Sub(maxY-1, maxX-1, 2, 0)
	separators := []int{
		0,  // the start
		30, // Name
		40, // Type
		60, // Cluster IP
		80, // Ports
	}

	w.Clear()
	w.ColorOn(colorHeader)
	w.HLine(0, 0, ' ', maxX)
	w.MovePrint(0, separators[0], " Name")
	w.MovePrint(0, separators[1], " Type")
	w.MovePrint(0, separators[2], " Cluster IP")
	w.MovePrint(0, separators[3], " Ports")
	w.ColorOff(colorHeader)

	t = &Tab{
		Panel:   panel,
		name:    "Services",
		entries: e,
		Redraw: func() error {
			subWindow.Clear()
			subMaxY, _ := subWindow.MaxYX()
			for i, entry := range e {
				if i >= subMaxY {
					return nil
				}

				subWindow.MovePrint(i, separators[0], " "+entry["NAME"])
				subWindow.MovePrint(i, separators[1], " "+entry["TYPE"])
				subWindow.MovePrint(i, separators[2], " "+entry["CLUSTER-IP"])
				subWindow.MovePrint(i, separators[3], " "+entry["PORTS"])
				i++
			}

			if err := w.Touch(); err != nil {
				return err
			}

			return nil
		},
		Update: func(o *gopenshift.OpenShift) error {
			services, err := o.GetServices()
			if err != nil {
				return err
			}

			e = []map[string]string{}
			for _, service := range services {
				portString := ""
				//for _, port := range service.Spec.Ports {
				//	portString = append(portString, strings.Join([]string{string(port.Port), string(port.Protocol)}, "/"))
				//	portString += ","
				//}

				e = append(e, map[string]string{
					"NAME":       service.Name,
					"TYPE":       string(service.Spec.Type),
					"CLUSTER-IP": service.Spec.ClusterIP,
					"PORTS":      portString,
				})
			}

			return nil
		},
	}

	return t
}
