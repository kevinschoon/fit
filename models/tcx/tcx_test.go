package tcx

import (
	"github.com/kevinschoon/gofit/models"
	"github.com/kevinschoon/tcx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTCX(t *testing.T) {
	data := &TCX{
		Acts: tcx.Acts{
			tcx.Activity{
				StartTime: time.Now(),
				Laps: []tcx.Lap{
					tcx.Lap{
						TotalTime: float64(time.Now().Sub(time.Now().AddDate(0, 0, 1)).Seconds()),
						Dist:      100.0,
					},
					tcx.Lap{
						TotalTime: float64(time.Now().Sub(time.Now().AddDate(0, 0, 1)).Seconds()),
						Dist:      100.0,
					},
				},
			},
			tcx.Activity{
				StartTime: time.Now().AddDate(0, 0, 1),
				Laps: []tcx.Lap{
					tcx.Lap{
						TotalTime: float64(time.Now().Sub(time.Now().AddDate(0, 0, 1)).Seconds()),
						Dist:      200.0,
					},
					tcx.Lap{
						TotalTime: float64(time.Now().Sub(time.Now().AddDate(0, 0, 1)).Seconds()),
						Dist:      200.0,
					},
				},
			},
		},
	}
	assert.Equal(t, 2, len(data.Acts))
	series := data.Load()[0]
	assert.Equal(t, 2, series.Len())
	assert.Equal(t, models.Value(200.0), series.Value(0, "Distance"))
	assert.Equal(t, models.Value(400.0), series.Value(1, "Distance"))
}

func TestTCXLoad(t *testing.T) {
	data, err := FromDir("test/sample.tcx")
	assert.NoError(t, err)
	series := data.Load()[0]
	assert.Equal(t, 1, series.Len())
	assert.Equal(t, models.Value(8348.5039063), series.Value(0, "Distance"))
	assert.Equal(t, models.Value(2325.02), series.Value(0, "Duration"))
	assert.Equal(t, models.Value(1.0), series.Value(0, "Laps"))
}
