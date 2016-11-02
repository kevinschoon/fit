package clients

import (
	"fmt"
	mtx "github.com/gonum/matrix/mat64"
	"github.com/kevinschoon/fit/types"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func DBPath() string {
	b := make([]rune, 5)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return fmt.Sprintf("/tmp/fit-test-", string(b))
}

func TestSQLiteClient(t *testing.T) {
	client, err := NewSQLClient(DBPath())
	assert.NoError(t, err)
	assert.NoError(client.Write(&types.Dataset{Name: "TestDS",
		Mtx: mtx.NewDense(2, 2, []float64{1.0, 2.0, 3.0, 4.0})}))
	result, err := client.Query(nil)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
