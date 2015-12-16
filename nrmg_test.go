package nrmg

import (
	"fmt"
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
	for {
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
		zones := AssignZones(diagram, zc)
		// for c, z := range zones {
		// 	log.Printf("cell %v zone %d", c.Site, z)
		// }

		log.Println(zc)
		unassignedSize, sizes := ZoneSizes(zc, zones, diagram)
		log.Printf("unassigned: %f |", unassignedSize)
		for ix, s := range sizes {
			fmt.Printf(" zone %d: %f |", ix, s)
		}

		loss := 0.
		for ix, s := range sizes {
			diff := s - zc.Sizes[ix]
			loss += diff * diff
		}

		log.Printf("loss: %f", loss)

		if loss < 0.0004 {
			RenderPreview(zones, diagram, "hello.png")
			break
		}
	}
}
