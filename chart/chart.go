package chart

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/kevinschoon/gofit/models"
	"image/color"
)

type PlotCfg struct {
	PrimaryColor   color.Color
	SecondaryColor color.Color
	XLabel         string
	YLabel         string
	Title          string
}

func getPlt(cfg PlotCfg) (*plot.Plot, error) {
	plt, err := plot.New()
	if err != nil {
		return nil, err
	}
	plt.Legend.Color = cfg.PrimaryColor
	plt.BackgroundColor = cfg.SecondaryColor

	plt.Title.Color = cfg.PrimaryColor
	plt.Title.Text = cfg.Title
	plt.Title.Font.Size = 0.5 * vg.Inch

	plt.Y.Color = cfg.PrimaryColor
	plt.Y.Label.Text = cfg.YLabel
	plt.Y.Label.Color = cfg.PrimaryColor
	plt.Y.Label.Font.Size = 0.3 * vg.Inch
	plt.Y.Tick.Color = cfg.PrimaryColor
	plt.Y.Tick.Label.Font.Size = 0.2 * vg.Inch
	plt.Y.Tick.Label.Color = cfg.PrimaryColor

	plt.X.Label.Text = cfg.XLabel
	plt.X.Color = cfg.PrimaryColor
	plt.X.Label.Color = cfg.PrimaryColor
	plt.X.Label.Font.Size = 0.3 * vg.Inch
	plt.X.Tick.Color = cfg.PrimaryColor
	plt.X.Tick.Label.Color = cfg.PrimaryColor
	plt.X.Tick.Label.Font.Size = 0.2 * vg.Inch
	//plt.X.Tick.Marker = plot.UnixTimeTicks{Format: "2006-01-02"}
	return plt, nil
}

func getXYs(collection *models.Collection, x, y models.Key) plotter.XYer {
	xys := make(plotter.XYs, collection.Len())
	for i, series := range collection.Series {
		xys[i].X = series.Sum(x)
		xys[i].Y = series.Sum(y)
	}
	return xys
}

// DistanceOverTime returns a chart displaying daily distance totals
func New(collection *models.Collection, x, y models.Key) (vg.CanvasWriterTo, error) {
	plt, _ := getPlt(PlotCfg{
		PrimaryColor:   color.White,
		SecondaryColor: color.Black,
		Title:          collection.Name,
		YLabel:         collection.GetName(y),
		XLabel:         collection.GetName(x),
	})
	line, _ := plotter.NewLine(getXYs(collection, x, y))
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
