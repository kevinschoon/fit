package main

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgsvg"
	"math/rand"
)

type ChartConfig struct {
	Title  string
	XLabel string
	YLabel string
	Filter Filter
}

type Chart interface {
	Canvas() (*vgsvg.Canvas, error)
}

type OverviewChart struct {
	Config *ChartConfig
}

func (chart OverviewChart) Canvas() (*vgsvg.Canvas, error) {
	p, err := newPlot(chart.Config)
	if err != nil {
		return nil, err
	}
	p.Add(plotter.NewGrid())
	laps := database.Laps(chart.Config.Filter)
	pts := make(plotter.Values, len(laps))
	for i := range pts {
		pts[i] = laps[i].Dist
	}
	barChart, err := plotter.NewBarChart(pts, vg.Points(20))
	if err != nil {
		return nil, err
	}
	p.Add(barChart)
	c := vgsvg.New(5*vg.Inch, 5*vg.Inch)
	// Draw to the Canvas.
	p.Draw(draw.New(c))
	return c, nil
}

type RegressionChart struct {
	Config *ChartConfig
}

func (chart RegressionChart) Canvas() (*vgsvg.Canvas, error) {
	// Get some data.
	n, m := 5, 10
	pts := make([]plotter.XYer, n)
	for i := range pts {
		xys := make(plotter.XYs, m)
		pts[i] = xys
		center := float64(i)
		for j := range xys {
			xys[j].X = center + (rand.Float64() - 0.5)
			xys[j].Y = center + (rand.Float64() - 0.5)
		}
	}

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}

	// Create two lines connecting points and error bars. For
	// the first, each point is the mean x and y value and the
	// error bars give the 95% confidence intervals.  For the
	// second, each point is the median x and y value with the
	// error bars showing the minimum and maximum values.
	mean95, err := plotutil.NewErrorPoints(plotutil.MeanAndConf95, pts...)
	if err != nil {
		panic(err)
	}
	medMinMax, err := plotutil.NewErrorPoints(plotutil.MedianAndMinMax, pts...)
	if err != nil {
		panic(err)
	}
	plotutil.AddLinePoints(plt,
		"mean and 95% confidence", mean95,
		"median and minimum and maximum", medMinMax)
	plotutil.AddErrorBars(plt, mean95, medMinMax)

	// Add the points that are summarized by the error points.
	plotutil.AddScatters(plt, pts[0], pts[1], pts[2], pts[3], pts[4])

	c := vgsvg.New(5*vg.Inch, 5*vg.Inch)
	// Draw to the Canvas.
	plt.Draw(draw.New(c))
	return c, nil

}

func newPlot(input *ChartConfig) (*plot.Plot, error) {
	p, err := plot.New()
	if err != nil {
		return nil, err
	}
	p.Title.Text = input.Title
	p.X.Label.Text = input.XLabel
	p.Y.Label.Text = input.YLabel
	return p, nil
}

func GetChart(name string, fn Filter) (chart Chart, err error) {
	switch name {
	case "regression":
		chart = RegressionChart{
			Config: &ChartConfig{
				Filter: fn,
			}}
	default:
		chart = OverviewChart{
			Config: &ChartConfig{
				Title:  "Overview",
				XLabel: "Days",
				YLabel: "Distance (Meters)",
				Filter: fn,
			}}
	}
	return chart, err
}
