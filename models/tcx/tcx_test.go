package tcx

import (
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
	collection := data.Load()
	assert.Equal(t, 2, collection.Len())
	assert.Equal(t, 200.0, collection.Series[0].Get(0, Distance).Value)
	assert.Equal(t, 400.0, collection.Series[1].Get(0, Distance).Value)
}

func TestTCXLoad(t *testing.T) {
	data, err := FromDir("test/sample.tcx")
	assert.NoError(t, err)
	collection := data.Load()
	assert.Equal(t, 1, collection.Len())
	assert.Equal(t, 8348.5039063, collection.Series[0].Get(0, Distance).Value)
	assert.Equal(t, 2325.02, collection.Series[0].Get(0, Duration).Value)
	assert.Equal(t, 1.0, collection.Series[0].Get(0, Laps).Value)
}
