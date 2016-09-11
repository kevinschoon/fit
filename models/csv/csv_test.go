package csv

import (
	"github.com/kevinschoon/gofit/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCSV(t *testing.T) {
	data, err := FromFile("test/LakeHuron.csv")
	assert.NoError(t, err)
	collection := data.Load()
	assert.Equal(t, 3, len(collection.Names()))
	assert.Equal(t, "", collection.Names()[0])
	assert.Equal(t, "time", collection.Names()[1])
	assert.Equal(t, "LakeHuron", collection.Names()[2])
	assert.Equal(t, 1, collection.Len())
	assert.Equal(t, 579.96, collection.Series[0].Get(97, models.Key(2)).Value)
}
