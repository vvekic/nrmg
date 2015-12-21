package nrmg

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/pzsz/voronoi/utils"
	"github.com/stretchr/testify/suite"
)

func init() {
	seed := time.Now().UnixNano()
	// seed = 1450220532465546300
	rand.Seed(seed)
	log.Printf("random seed: %d", seed)
}

type ZoneCreation struct {
	suite.Suite
}

func TestZoneCreation(t *testing.T) {
	suite.Run(t, new(ZoneCreation))
}

func (s *ZoneCreation) TestVoronoi() {

	diagram := Voronoi(200)
	totalArea := 0.
	for _, c := range diagram.Cells {
		totalArea -= utils.CellArea(c)
	}
	log.Println(totalArea)

	zc := &ZoneConfig{
		Sizes: []float64{25, 12, 13, 13, 12, 25},
	}
	zc.NormalizeSizes()
	m := New()
	m.AssignZones(diagram, zc)
	m.SaveImage("map.png")
}
