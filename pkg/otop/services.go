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

	t = &Tab{
		Panel:   panel,
		name:    "Services",
		entries: e,
		Initialize: func() error {
			w.Clear()
			w.ColorOn(colorHeader)
			w.HLine(0, 0, ' ', maxX)
			w.MovePrint(0, 0, " Name")
			w.MovePrint(0, 15, " Type")
			w.MovePrint(0, 20, " Cluster IP")
			w.MovePrint(0, 30, " Ports")
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
				subWindow.MovePrint(i, 15, " "+entry["TYPE"])
				subWindow.MovePrint(i, 20, " "+entry["CLUSTER-IP"])
				subWindow.MovePrint(i, 30, " "+entry["PORTS"])
				i++
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
