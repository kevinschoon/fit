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

func Values(totals Totals) plotter.Values {
	pts := make(plotter.Values, len(totals))
	for i := range pts {
		pts[i] = totals[i].Km()
	}
	return pts
}

func Names(totals Totals) []string {
	names := make([]string, len(totals))
	for i := range names {
		names[i] = totals[i].Name()
	}
	return names
}

type Chart interface {
	Canvas(Totals) (*vgsvg.Canvas, error)
}

type OverviewChart struct {
	Title string
}

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
	barChart, err := plotter.NewBarChart(Values(totals), vg.Points(10))
	if err != nil {
		return nil, err
	}
	plt.Add(barChart)
	canvas := vgsvg.New(5*vg.Inch, 5*vg.Inch)
	plt.Draw(draw.New(canvas))
	return canvas, nil
}

type RegressionChart struct{}

func (chart RegressionChart) Canvas(totals Totals) (*vgsvg.Canvas, error) {
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
	p.Title.Text = "Horizontal box plots"
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
