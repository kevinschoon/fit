package chart

import (
	mtx "github.com/gonum/matrix/mat64"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"image/color"
)

type Config struct {
	PrimaryColor   color.Color
	SecondaryColor color.Color
	XLabel         string
	YLabel         string
	Title          string
	Type           string
	Columns        []string
	Width          vg.Length
	Height         vg.Length
	PlotTime       bool
}

type Vector struct {
	Name     string
	Position int
}

func getPlot(cfg Config) (*plot.Plot, error) {
	plt, err := plot.New()
	if err != nil {
		return nil, err
	}
	plt.Legend.Color = cfg.PrimaryColor
	plt.Legend.Top = true
	plt.BackgroundColor = cfg.SecondaryColor

	//plt.Title.Text = cfg.Title
	//plt.Title.Color = cfg.PrimaryColor
	//plt.Title.Font.Size = 0.5 * vg.Inch

	//plt.Y.Label.Text = cfg.YLabel
	plt.Y.Color = cfg.PrimaryColor
	plt.Y.Label.Color = cfg.PrimaryColor
	plt.Y.Label.Font.Size = 0.3 * vg.Inch
	plt.Y.Tick.Color = cfg.PrimaryColor
	plt.Y.Tick.Label.Font.Size = 0.2 * vg.Inch
	plt.Y.Tick.Label.Color = cfg.PrimaryColor

	//plt.X.Label.Text = cfg.XLabel
	plt.X.Color = cfg.PrimaryColor
	plt.X.Label.Color = cfg.PrimaryColor
	plt.X.Label.Font.Size = 0.3 * vg.Inch
	plt.X.Tick.Color = cfg.PrimaryColor
	plt.X.Tick.Label.Color = cfg.PrimaryColor
	plt.X.Tick.Label.Font.Size = 0.2 * vg.Inch

	if cfg.PlotTime {
		plt.X.Tick.Marker = plot.UnixTimeTicks{Format: "2006-01-02"}
	}
	return plt, nil
}

// GetLines returns variadic arguments for plotter.AddLines
// The chart will always plot the first column vector as
// the X axis and remaining columns on the Y axis.
//
//	x,y,z
//  1,1,2
//  2,3,4
//  3,5,6
//
// Line Y: 1,1 2,3 3,5
// Line Z: 1,2 2,4 3,6

func GetLines(mx *mtx.Dense, columns []string) []interface{} {
	data := make([]interface{}, 0)
	r, c := mx.Dims()
	for i := 1; i < c; i++ {
		xys := make(plotter.XYs, r)
		for j := 0; j < r; j++ {
			xys[j].X = mx.At(j, 0)
			xys[j].Y = mx.At(j, i)
		}
		data = append(data, columns[i])
		data = append(data, xys)
	}
	return data
}

// GetValues returns an array of plotter.Values
// where each entry is a column vector.
func GetValues(mx *mtx.Dense) []plotter.Values {
	_, c := mx.Dims()
	values := make([]plotter.Values, c)
	for i := 0; i < c; i++ {
		values[i] = plotter.Values(mtx.Col(nil, i, mx))
	}
	return values
}

func New(cfg Config, mx *mtx.Dense) (vg.CanvasWriterTo, error) {
	plt, err := getPlot(cfg)
	if err != nil {
		return nil, err
	}
	switch cfg.Type {
	case "box":
		values := GetValues(mx)
		for i, vals := range values {
			box, err := plotter.NewBoxPlot(vg.Points(20), float64(i), vals)
			if err != nil {
				return nil, err
			}
			box.WhiskerStyle.Color = cfg.PrimaryColor
			box.BoxStyle.Color = cfg.PrimaryColor
			box.MedianStyle.Color = cfg.PrimaryColor
			plt.Add(box)
		}
	default: // Default to line chart
		if err := plotutil.AddLines(plt, GetLines(mx, cfg.Columns)...); err != nil {
			return nil, err
		}
	}
	plt.Add(plotter.NewGrid())
	canvas, err := draw.NewFormattedCanvas(cfg.Width, cfg.Height, "png")
	if err != nil {
		return nil, err
	}
	plt.Draw(draw.New(canvas))
	return canvas, nil
}
