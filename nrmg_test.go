package nrmg

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/persomi/set"
	"github.com/pzsz/voronoi"
	"github.com/stretchr/testify/assert"
)

func init() {
	seed := time.Now().UnixNano()
	// seed = 1450220532465546300
	rand.Seed(seed)
	log.Printf("random seed: %d", seed)
}

func newTestDiagram() *voronoi.Diagram {
	return createDiagram(10)
}

func newTestMap() *Map {
	m := New()
	return m
}
func TestAddCells(t *testing.T) {
	td := newTestDiagram()
	cells := td.Cells
	cs := newCellSet()
	addCells(cells, cs)
	assert.Equal(t, len(cells), cs.(set.Interface).Size())
	for _, c := range cells {
		assert.True(t, cs.(set.Interface).Has(c))
	}
}

func TestRandomChoice(t *testing.T) {
	td := newTestDiagram()
	cs := newCellSet()
	addCells(td.Cells, cs)
	rc := randomChoice(cs)
	assert.True(t, cs.(set.Interface).Has(rc))
}

func TestOtherZones(t *testing.T) {
	m := newTestMap()
	m.zoneCount = 5
	oz := m.otherZones(2)
	assert.Equal(t, []int{0, 1, 3, 4}, oz)
}

func TestNeighbors(t *testing.T) {
	m := newTestMap()
	m.d = createDiagram(10)
	cs := newCellSet()
	addCells(m.d.Cells, cs)
	for i := 0; i < 10; i++ {
		rc := randomChoice(cs)
		n := m.neighbors(cs.List().([]*voronoi.Cell)...)

	}
}

func TestTesselate(t *testing.T) {
	t.SkipNow()
	c := &Config{
		ZoneSizes: []float64{25, 12, 13, 13, 12, 25},
		Cells:     200,
	}
	c.NormalizeSizes()
	log.Println(c)
	m := New()
	m.tesselate(c)
	m.saveImage("map.png")
}
