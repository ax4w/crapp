package tools

import (
	"crapp/internal/bridge"
	"fmt"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func makeLabels(data []bridge.HRObj) []string {
	labels := make([]string, len(data))
	for i, v := range data {
		labels[i] = fmt.Sprintf("%d bpm", v.HR)
	}
	return labels
}

func ShowGraph(data []bridge.HRObj) {
	if len(data) == 0 {
		return
	}
	p := plot.New()
	p.Title.Text = "HR Vis"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "HR"
	p.X.Tick.Marker = plot.TimeTicks{
		Format: "15:04:05", // Change this format as needed
		Time:   func(t float64) time.Time { return time.Unix(int64(t), 0) },
	}

	var pts plotter.XYs
	for i, v := range data {
		pts = append(pts, plotter.XY{
			X: float64(v.T.Unix()),
			Y: float64(v.HR),
		})
		println("point", i, "x", pts[i].X, "y", pts[i].Y)
	}
	err := plotutil.AddLinePoints(p, "Data", pts)
	labels, _ := plotter.NewLabels(plotter.XYLabels{
		XYs:    pts,
		Labels: makeLabels(data),
	})
	labels.Offset.X = vg.Points(0)
	labels.Offset.Y = vg.Points(5)
	p.Add(labels)
	if err != nil {
		println(err.Error())
		return
	}
	if err := p.Save(30*vg.Inch, 10*vg.Inch, "hr.svg"); err != nil {
		panic(err)
	}
}
