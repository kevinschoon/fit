package types

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	u, err := url.Parse("http://localhost/?q=D0,x,y,z&q=D1,z&grouping=Duration,0,1m&fn=avg")
	assert.NoError(t, err)
	query := NewQueryQS(u)
	assert.Equal(t, 2, query.Len())
	assert.Equal(t, "D0", query.Datasets[0].Name)
	assert.Equal(t, 3, len(query.Datasets[0].Columns))
	assert.Equal(t, "x", query.Datasets[0].Columns[0])
	assert.Equal(t, "y", query.Datasets[0].Columns[1])
	assert.Equal(t, "z", query.Datasets[0].Columns[2])
	assert.Equal(t, 1, len(query.Datasets[1].Columns))
	assert.Equal(t, "z", query.Datasets[1].Columns[0])
	assert.Equal(t, "D1", query.Datasets[1].Name)
	assert.Equal(t, "avg", query.Function.Name)
	assert.Equal(t, 0, query.Grouping.Index)
	assert.Equal(t, time.Minute, query.Grouping.Max)
	assert.Equal(t, "fn=avg&grouping=Duration%2C0%2C1m0s&q=D0%2Cx%2Cy%2Cz&q=D1%2Cz", query.String())
}
