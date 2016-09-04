package main

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"image/color"
)

// DistanceOverTime returns a chart displaying daily distance totals
func DistanceOverTime(data plotter.XYer) (vg.CanvasWriterTo, error) {
	plt, err := plot.New()
	if err != nil {
		return nil, err
	}
	plt.Legend.Color = color.White
	plt.BackgroundColor = color.Black

	plt.Title.Color = color.White
	plt.Title.Text = "Distance Traveled"
	plt.Title.Font.Size = 0.5 * vg.Inch

	plt.Y.Color = color.White
	plt.Y.Label.Text = "Distance (Km)"
	plt.Y.Label.Color = color.White
	plt.Y.Label.Font.Size = 0.3 * vg.Inch
	plt.Y.Tick.Color = color.White
	plt.Y.Tick.Label.Font.Size = 0.2 * vg.Inch
	plt.Y.Tick.Label.Color = color.White

	plt.X.Label.Text = "Time"
	plt.X.Color = color.White
	plt.X.Label.Color = color.White
	plt.X.Label.Font.Size = 0.3 * vg.Inch
	plt.X.Tick.Color = color.White
	plt.X.Tick.Label.Color = color.White
	plt.X.Tick.Label.Font.Size = 0.2 * vg.Inch
	plt.X.Tick.Marker = plot.UnixTimeTicks{Format: "2006-01-02"}

	line, err := plotter.NewLine(data)
	if err != nil {
		panic(err)
	}
	line.LineStyle.Color = color.White
	line.LineStyle.Width = vg.Points(2)

	plt.Add(plotter.NewGrid())
	plt.Add(line)

	canvas, err := draw.NewFormattedCanvas(32*vg.Inch, 12*vg.Inch, "png")
	if err != nil {
		panic(err)
	}
	plt.Draw(draw.New(canvas))
	return canvas, nil
}
