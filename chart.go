package main

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgsvg"
	"time"
)

// DistanceOverTime returns a chart displaying daily distance totals
func DistanceOverTime(data map[time.Time]float64) (*vgsvg.Canvas, error) {
	plt, err := plot.New()
	if err != nil {
		return nil, err
	}
	plt.Title.Text = ""
	plt.Y.Label.Text = "Y"
	plt.X.Label.Text = "X"
	plt.Add(plotter.NewGrid())
	l, err := plotter.NewLine(lineData)
	if err != nil {
		panic(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}
	canvas := vgsvg.New(5*vg.Inch, 5*vg.Inch)
	plt.Draw(draw.New(canvas))
	return canvas, nil
}

/*
func (chart OverviewChart) Canvas(totals Totals) (*vgsvg.Canvas, error) {
	plt, err := plot.New()
	if err != nil {
		return nil, err
	}
	plt.Title.Text = chart.Title
	plt.Y.Label.Text = "Km"
	plt.X.Label.Text = "Time Period"
	plt.NominalX(Names(totals)...)
	plt.Add(plotter.NewGrid())
	barChart, err := plotter.NewBarChart(Values(totals), vg.Points(1))
	if err != nil {
		return nil, err
	}
	plt.Add(barChart)
	canvas := vgsvg.New(5*vg.Inch, 5*vg.Inch)
	plt.Draw(draw.New(canvas))
	return canvas, nil
}
*/

/*
type RegressionChart struct{}

func (chart RegressionChart) Canvas(totals Totals) (*vgsvg.Canvas, error) {
	plt, err := plot.New()
	if err != nil {
		return nil, err
	}
	pts := make(plotter.XYs, len(totals))
	for i := range pts {
		pts[i].X = totals[i].TotalTime
		pts[i].Y = totals[i].Dist
	}
	s, err := plotter.NewScatter(pts)
	if err != nil {
		return nil, err
	}
	plt.Add(s)
	l, err := plotter.NewLine(SeriesToXYs(totals.Predict()))
	if err != nil {
		return nil, err
	}
	plt.Add(l)
	c := vgsvg.New(5*vg.Inch, 5*vg.Inch)
	// Draw to the Canvas.
	plt.Draw(draw.New(c))
	return c, nil
}

type DistributionChart struct{}

func (chart DistributionChart) Canvas(totals Totals) (*vgsvg.Canvas, error) {
	// Get some data to display in our plot.
	rand.Seed(int64(0))
	n := 10
	uniform := make(plotter.Values, n)
	normal := make(plotter.Values, n)
	expon := make(plotter.Values, n)
	for i := 0; i < n; i++ {
		uniform[i] = rand.Float64()
		normal[i] = rand.NormFloat64()
		expon[i] = rand.ExpFloat64()
	}

	// Create the plot and set its title and axis label.
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Distribution"
	p.X.Label.Text = "Values"

	// Make horizontal boxes for our data and add
	// them to the plot.
	w := vg.Points(20)
	b0, err := plotter.NewBoxPlot(w, 0, uniform)
	if err != nil {
		panic(err)
	}
	b1, err := plotter.NewBoxPlot(w, 1, normal)
	if err != nil {
		panic(err)
	}
	b2, err := plotter.NewBoxPlot(w, 2, expon)
	if err != nil {
		panic(err)
	}
	p.Add(b0, b1, b2)

	// Set the Y axis of the plot to nominal with
	// the given names for y=0, y=1 and y=2.
	p.NominalY("Uniform\nDistribution", "Normal\nDistribution",
		"Exponential\nDistribution")

	c := vgsvg.New(5*vg.Inch, 5*vg.Inch)
	// Draw to the Canvas.
	p.Draw(draw.New(c))
	return c, nil
}
*/
