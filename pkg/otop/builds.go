package otop

import (
	"fmt"
	"regexp"

	"github.com/oatmealraisin/gopenshift/pkg/gopenshift"
	buildapi "github.com/openshift/origin/pkg/build/api"
	gc "github.com/rthornton128/goncurses"
)

var removeJenkinsInfo = regexp.MustCompile(":.*")

func NewBuildsTab(w *gc.Window) *Tab {
	panel := gc.NewPanel(w)
	var t *Tab
	maxY, maxX := w.MaxYX()
	subWindow := w.Sub(maxY-1, maxX-1, 2, 0)

	builds := []*buildapi.Build{}

	separators := []int{
		0,  // the start
		8,  // Type (Needs to be this to cut off Pipeline in JenkinsPipeline)
		28, // Triggered by
	}

	w.Clear()
	w.ColorOn(colorHeader)
	w.HLine(0, 0, ' ', maxX)
	w.MovePrint(0, separators[0], " Type")
	w.MovePrint(0, separators[1], " Triggered By")
	w.MovePrint(0, separators[len(separators)-1], " Name")
	w.ColorOff(colorHeader)

	t = &Tab{
		Panel: panel,
		name:  "Builds",
		Redraw: func() error {
			subWindow.Clear()
			subMaxY, _ := subWindow.MaxYX()
			for i, build := range builds {
				if i >= subMaxY {
					return nil
				}

				subWindow.MovePrint(i, separators[0], fmt.Sprintf(" %.8s", buildapi.StrategyType(build.Spec.Strategy)))
				subWindow.MovePrint(i, separators[1], fmt.Sprintf(" %.20s", removeJenkinsInfo.ReplaceAllString(build.Spec.TriggeredBy[0].Message, "")))
				subWindow.MovePrint(i, separators[len(separators)-1], fmt.Sprintf(" %s", build.Name))
			}

			if err := w.Touch(); err != nil {
				return err
			}

			return nil
		},
		Update: func(o *gopenshift.OpenShift) error {
			var err error
			builds, err = o.GetBuilds()
			return err
		},
	}

	return t
}
