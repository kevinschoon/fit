package chart

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/kevinschoon/gofit/models"
	"image/color"
)

type Config struct {
	PrimaryColor   color.Color
	SecondaryColor color.Color
	XLabel         string
	YLabel         string
	Title          string
	Keys           []string
	XAxis          models.Key
	YAxis          map[string]models.Key
	Width          vg.Length
	Height         vg.Length
}

func getPlot(cfg Config) (*plot.Plot, error) {
	plt, err := plot.New()
	if err != nil {
		return nil, err
	}
	plt.Legend.Color = cfg.PrimaryColor
	plt.Legend.Top = true
	plt.Legend.YOffs = 0.1 * vg.Inch
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
	plt.X.Tick.Marker = plot.UnixTimeTicks{Format: "2006-01-02"}
	return plt, nil
}

// New returns a line chart plotting the provided values
func New(cfg Config, series []*models.Series) (vg.CanvasWriterTo, error) {
	plt, err := getPlot(cfg)
	if err != nil {
		return nil, err
	}
	// Implement variadic arguments for plotter.AddLines
	// ...string,XYer
	data := make([]interface{}, 0)
	xValues := models.Select(cfg.XAxis, series)
	for _, name := range models.Keys(series[0]) {
		if y, ok := cfg.YAxis[name]; ok {
			xys := make(plotter.XYs, len(xValues))
			for i, value := range models.Select(y, series) {
				xys[i].X = xValues[i].Float64()
				xys[i].Y = value.Float64()
			}
			data = append(data, name)
			data = append(data, xys)
		}
	}
	if err := plotutil.AddLines(plt, data...); err != nil {
		return nil, err
	}
	plt.Add(plotter.NewGrid())
	canvas, err := draw.NewFormattedCanvas(cfg.Width, cfg.Height, "png")
	if err != nil {
		return nil, err
	}
	plt.Draw(draw.New(canvas))
	return canvas, nil
}
