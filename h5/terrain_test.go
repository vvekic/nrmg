package h5

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testTerrainFile = "fixture/GroundTerrain.bin"
)

func TestTerrain(t *testing.T) {
	td, err := NewTerrainData(testTerrainFile)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.NotNil(t, td) {
		t.FailNow()
	}
	err = td.Save(testTerrainFile + ".saved")
	assert.NoError(t, err)
}
