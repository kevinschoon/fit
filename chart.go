package main

import (
	"fmt"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgsvg"
	"github.com/montanaflynn/stats"
	"image/color"
)

func RegressXYs(data plotter.XYer) plotter.XYer {
	series := make(stats.Series, data.Len())
	regressions := make(plotter.XYs, data.Len())
	for i := range series {
		series[i].X, series[i].Y = data.XY(i)
	}
	r, err := stats.LinearRegression(series)
	if err != nil {
		panic(err)
	}
	for i, k := range r {
		regressions[i].X = k.X
		regressions[i].Y = k.Y
		fmt.Println(series[i].X, series[i].Y, "--> ", k.X, k.Y)
	}
	return regressions
}

// DistanceOverTime returns a chart displaying daily distance totals
func DistanceOverTime(data plotter.XYer) (*vgsvg.Canvas, error) {
	plt, err := plot.New()
	if err != nil {
		return nil, err
	}
	//bars, err := plotter.NewBarChart(data, vg.Points(1))
	plt.Title.Text = "Distance Traveled"
	plt.Y.Label.Text = "Distance (Km)"
	plt.X.Label.Text = "Time"
	plt.X.Tick.Marker = plot.UnixTimeTicks{Format: "2006-01-02"}
	lines, err := plotter.NewLine(data)
	if err != nil {
		panic(err)
	}
	plt.Add(plotter.NewGrid())
	//plt.Add(scatter)
	//plt.Add(bars)
	plt.Add(lines)
	canvas := vgsvg.New(20*vg.Inch, 4*vg.Inch)
	plt.Draw(draw.New(canvas))
	return canvas, nil
}

func ChartXYs(data plotter.XYer) (*vgsvg.Canvas, error) {
	plt, err := plot.New()
	if err != nil {
		return nil, err
	}
	plt.Title.Text = ""
	plt.Y.Label.Text = "Distance (km)"
	plt.X.Label.Text = ""
	plt.X.Tick.Marker = plot.UnixTimeTicks{Format: "2006-01-02"}
	l, err := plotter.NewLine(data)
	if err != nil {
		return nil, err
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}
	plt.Add(plotter.NewGrid())
	plt.Add(l)
	canvas := vgsvg.New(20*vg.Inch, 4*vg.Inch)
	plt.Draw(draw.New(canvas))
	return canvas, nil
}
