package nrmg

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/pzsz/voronoi"
	"github.com/pzsz/voronoi/utils"
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
	cs := voronoi.NewCellSet(cells...)
	assert.Equal(t, len(cells), cs.Cardinality())
	assert.True(t, cs.ContainsAll(cells...))
}

func TestRandomChoice(t *testing.T) {
	td := newTestDiagram()
	cs := voronoi.NewCellSet(td.Cells...)
	rc := randomChoice(cs)
	assert.True(t, cs.Contains(rc))
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
	cs := voronoi.NewCellSet()
	addCells(m.d.Cells, cs)
	for i := 0; i < 10; i++ {
		rc := randomChoice(cs)
		rc2 := randomChoice(cs.Difference(voronoi.NewCellSet(rc)))
		n := m.neighbors(voronoi.NewCellSet(rc, rc2), false)
		assert.False(t, n.Contains(rc))
		assert.False(t, n.Contains(rc2))
		for c := range n.Iter() {
			assert.True(t, utils.AreNeighbors(c, rc) || utils.AreNeighbors(c, rc2))
		}
	}
}

func TestAssign(t *testing.T) {
	m := newTestMap()
	m.d = createDiagram(10)
	firstCell := m.d.Cells[0]
	err := m.assign(0, firstCell)
	assert.NoError(t, err)
	assert.Equal(t, 1, m.zoneCell[0].Cardinality())
	assert.True(t, m.zoneCell[0].Contains(firstCell))
	err = m.assign(0, firstCell)
	assert.Error(t, err)
}

func TestAssignSet(t *testing.T) {
	m := newTestMap()
	m.d = createDiagram(10)
	cs := voronoi.NewCellSet(m.d.Cells[0:5]...)
	err := m.assignSet(0, cs)
	assert.NoError(t, err)
	zc := m.zoneCells(0)
	assert.True(t, zc.Difference(cs).IsEmpty())
	err = m.assignSet(0, cs)
	assert.Error(t, err)
}

func TestZoneNeighbors(t *testing.T) {
	m := newTestMap()
	m.d = createDiagram(10)
	c := voronoi.NewCellSet(m.d.Cells...).Iter()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			m.assign(i, <-c)
		}
	}
	n := m.zoneNeighbors(0, 2)
	assert.True(t, n.Intersect(m.zoneCells(0)).IsEmpty())
	assert.True(t, n.Intersect(m.zoneCells(1)).IsEmpty())
	assert.True(t, n.Intersect(m.zoneCells(2)).IsEmpty())
}

func TestZoneCells(t *testing.T) {
	m := newTestMap()
	m.d = createDiagram(10)
	c := m.d.Cells[0]
	m.assign(0, c)
	zc0 := m.zoneCells(0)
	assert.Equal(t, 1, zc0.Cardinality())
	assert.True(t, zc0.Contains(c))
	zc1 := m.zoneCells(1)
	assert.True(t, zc1.IsEmpty())
}

func TestGrowZone(t *testing.T) {
	m := newTestMap()
	m.d = createDiagram(10)
	m.c = &Config{ZoneSizes: []float64{0.5, 0.8}}
	startCell := m.d.Cells[0]
	grown, err := m.growZone(0, startCell)
	assert.NoError(t, err)
	// log.Printf("cellsetarea: %f", utils.CellSetArea(grown))
	assert.True(t, utils.CellSetArea(grown) >= 0.5)
	m.assignSet(0, grown)
	startCell = randomChoice(grown)
	grown, err = m.growZone(1, startCell)
	assert.Error(t, err)
}

func TestTesselate(t *testing.T) {
	// t.SkipNow()
	c := &Config{
		ZoneSizes: []float64{25, 12, 13, 13, 12, 25},
		Cells:     200,
	}
	c.NormalizeSizes()
	log.Println(c)
	for {
		m := New()
		m.tesselate(c)
		e := m.zoneError()
		if e < 0.005 {
			log.Printf("sq err: %f", e)
			m.saveImage("map.png")
			break
		}
	}
}
