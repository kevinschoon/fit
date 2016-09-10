package chart

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/kevinschoon/gofit/models"
	"image/color"
)

func getPlt(primary, secondary color.Color) (*plot.Plot, error) {
	plt, err := plot.New()
	if err != nil {
		return nil, err
	}
	plt.Legend.Color = primary
	plt.BackgroundColor = secondary

	plt.Title.Color = primary
	plt.Title.Text = "TITLE"
	plt.Title.Font.Size = 0.5 * vg.Inch

	plt.Y.Color = primary
	plt.Y.Label.Text = "Y LABEL"
	plt.Y.Label.Color = primary
	plt.Y.Label.Font.Size = 0.3 * vg.Inch
	plt.Y.Tick.Color = primary
	plt.Y.Tick.Label.Font.Size = 0.2 * vg.Inch
	plt.Y.Tick.Label.Color = primary

	plt.X.Label.Text = "X LABEL"
	plt.X.Color = primary
	plt.X.Label.Color = primary
	plt.X.Label.Font.Size = 0.3 * vg.Inch
	plt.X.Tick.Color = primary
	plt.X.Tick.Label.Color = primary
	plt.X.Tick.Label.Font.Size = 0.2 * vg.Inch
	plt.X.Tick.Marker = plot.UnixTimeTicks{Format: "2006-01-02"}
	return plt, nil
}

func getXYs(collection *models.Collection) plotter.XYer {
	xys := make(plotter.XYs, collection.Len())
	for i, series := range collection.Series {
		xys[i].X = series.Sum(models.Key(0))
		xys[i].Y = series.Sum(models.Key(1))
	}
	return xys
}

// DistanceOverTime returns a chart displaying daily distance totals
func New(collection *models.Collection) (vg.CanvasWriterTo, error) {
	plt, err := getPlt(color.White, color.Black)
	line, err := plotter.NewLine(getXYs(collection))
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
