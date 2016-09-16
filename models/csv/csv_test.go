package csv

import (
	"github.com/kevinschoon/gofit/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCSV(t *testing.T) {
	data, err := FromFile("test/LakeHuron.csv")
	assert.NoError(t, err)
	series := data.Load()
	assert.Equal(t, 1, len(series))
	assert.Equal(t, 4, len(series[0].Keys))
	assert.Equal(t, models.Key(0), series[0].Keys["time"])
	assert.Equal(t, models.Key(1), series[0].Keys["NO_NAME"])
	assert.Equal(t, models.Key(2), series[0].Keys["_time"])
	assert.Equal(t, models.Key(3), series[0].Keys["LakeHuron"])
	assert.Equal(t, models.Value(579.96), series[0].Value(97, "LakeHuron"))
}
