package nrmg

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"sort"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/pzsz/voronoi"
	"github.com/pzsz/voronoi/utils"
)

type Config struct {
	ZoneSizes []float64
	Cells     int
}

func (c *Config) NormalizeSizes() {
	total := 0.
	for _, s := range c.ZoneSizes {
		total += s
	}
	for ix, s := range c.ZoneSizes {
		c.ZoneSizes[ix] = s / total
	}
}

func (c Config) String() string {
	var out string
	for ix, s := range c.ZoneSizes {
		out += fmt.Sprintf(" zone %d: %f |", ix, s)
	}
	return out
}

type Map struct {
	d         *voronoi.Diagram
	cellZone  map[*voronoi.Cell]int
	zoneCell  map[int]*voronoi.CellSet
	c         *Config
	zoneCount int
}

func New() *Map {
	return &Map{
		cellZone: map[*voronoi.Cell]int{},
		zoneCell: map[int]*voronoi.CellSet{},
	}
}

func (m *Map) clone() *Map {
	cl := New()
	cl.c = m.c
	cl.d = m.d
	for c, z := range m.cellZone {
		cl.cellZone[c] = z
	}
	for z, cs := range m.zoneCell {
		cl.zoneCell[z] = cs.Clone()
	}
	return cl
}

func (m *Map) tesselate(c *Config) {
	m.d = createDiagram(c.Cells)
	m.c = c
	m.zoneCount = len(c.ZoneSizes)
	grow := voronoi.NewCellSet()
	addCells(m.d.Cells, grow)
	for zoneId := 0; zoneId < m.zoneCount; zoneId++ {
		grown, err := m.tryGrowingZone(zoneId, grow)
		if err != nil {
			grow = m.zoneNeighbors(m.otherZones(zoneId)...)
			grow = grow.Difference(grown)
			grow = grow.Difference(m.zoneNeighbors(zoneId))
			grown, err = m.tryGrowingZone(zoneId, grow)
			if err != nil {
				log.Printf("Zone %d error: %v", zoneId, err)
				m.assignSet(zoneId, grown)
				continue
			}
		}
		m.assignSet(zoneId, grown)
		grow = m.neighbors(grown, false)
	}
}

func createDiagram(n int) *voronoi.Diagram {
	bb := voronoi.NewBBox(0., 1., 0., 1.)
	sites := utils.RandomSites(bb, n)
	d := voronoi.ComputeDiagram(sites, bb, true)
	for i := 0; i < 1; i++ {
		sites = utils.LloydRelaxation(d.Cells)
		d = voronoi.ComputeDiagram(sites, bb, true)
	}
	return d
}

func (m *Map) tryGrowingZone(zoneId int, grow *voronoi.CellSet) (*voronoi.CellSet, error) {
	bad := voronoi.NewCellSet()
	best := voronoi.NewCellSet()
	for {
		if grow.IsEmpty() {
			return best, fmt.Errorf("error zone %d: grow set empty", zoneId)
		}
		origin := randomChoice(grow)
		grown, err := m.growZone(zoneId, origin)
		if err != nil {
			if utils.CellSetArea(grown) > utils.CellSetArea(best) {
				best = grown.Clone()
			}
			bad = bad.Union(grown)
			grow = grow.Difference(bad)
			if grow.IsEmpty() {
				return best, err
			}
			continue
		}
		return grown, nil
	}
}

func (m *Map) growZone(zoneId int, startCell *voronoi.Cell) (*voronoi.CellSet, error) {
	cl := m.clone()
	area := 0.
	prescribedArea := cl.c.ZoneSizes[zoneId]
	grow := voronoi.NewCellSet()
	grow.Add(startCell)
	grown := voronoi.NewCellSet()
	for {
		// log.Printf("area: %f", area)
		if area > prescribedArea {
			return grown, nil
		}
		// log.Printf("grow: %d", grow.Cardinality())
		if grow.IsEmpty() {
			return grown, fmt.Errorf("No space for zone %d", zoneId)
		}

		var next *voronoi.Cell
		if cl.zoneCells(zoneId).Cardinality() == 0 {
			next = randomChoice(grow)
		} else {
			next = weightedChoice(cl.setWeights(grow, zoneId))
		}
		grow.Remove(next)
		// log.Println(grow)
		cl.assign(zoneId, next)
		grown.Add(next)
		// log.Printf("d_area: %f", utils.CellArea(next))
		area -= utils.CellArea(next)
		n := cl.neighbors(voronoi.NewCellSet(next), true)
		grow = grow.Union(n)
		// cl.saveImage(fmt.Sprintf("temp/%d.png", time.Now().UnixNano()))
	}
}

func (m *Map) assignSet(zoneId int, cs *voronoi.CellSet) error {
	for c := range cs.Iter() {
		if err := m.assign(zoneId, c); err != nil {
			return err
		}
	}
	return nil
}

func (m *Map) assign(zoneId int, cells ...*voronoi.Cell) error {
	for _, c := range cells {
		if _, ok := m.cellZone[c]; !ok {
			m.cellZone[c] = zoneId
			if m.zoneCell[zoneId] == nil {
				m.zoneCell[zoneId] = voronoi.NewCellSet()
			}
			m.zoneCell[zoneId].Add(c)
		} else {
			return fmt.Errorf("Cell already assigned: site: %v, zone: %d", c.Site, zoneId)
		}
	}
	return nil
}

func (m *Map) zoneCells(zoneId int) *voronoi.CellSet {
	if c, ok := m.zoneCell[zoneId]; !ok {
		return voronoi.NewCellSet()
	} else {
		return c
	}
}

func (m *Map) zoneNeighbors(zoneIds ...int) *voronoi.CellSet {
	out := voronoi.NewCellSet()
	for _, zoneId := range zoneIds {
		cs := m.zoneCells(zoneId)
		out = out.Union(m.neighbors(cs, false))
	}
	return out
}

func (m *Map) neighbors(cs *voronoi.CellSet, thickOnly bool) *voronoi.CellSet {
	out := voronoi.NewCellSet()
	for c := range cs.Iter() {
		ocSet := voronoi.NewCellSet()
		for _, he := range c.Halfedges {
			if utils.Distance(he.GetStartpoint(), he.GetEndpoint()) < 0.03 {
				continue
			}
			oc := he.Edge.GetOtherCell(c)
			if oc != nil && !cs.Contains(oc) {
				if _, ok := m.cellZone[oc]; !ok {
					ocSet.Add(oc)
				}
			}
		}
		out = out.Union(ocSet)
	}
	return out
}

func (m *Map) otherZones(zoneId int) []int {
	var out []int
	for i := 0; i < m.zoneCount; i++ {
		if i != zoneId {
			out = append(out, i)
		}
	}
	return out
}

func randomChoice(cs *voronoi.CellSet) *voronoi.Cell {
	list := cs.ToSlice()
	return list[rand.Intn(len(list))]
}

func weightedChoice(ws map[*voronoi.Cell]float64) *voronoi.Cell {
	cdf := []float64{}
	index := map[int]*voronoi.Cell{}
	cumulative := 0.
	i := 0
	for c, w := range ws {
		cdf = append(cdf, cumulative+w)
		index[i] = c
		i++
	}
	r := rand.Float64() * cdf[len(cdf)-1]
	chosenIx := sort.Search(len(cdf), func(i int) bool { return cdf[i] >= r })
	return index[chosenIx]
}

func (m *Map) setWeights(cs *voronoi.CellSet, zoneId int) map[*voronoi.Cell]float64 {
	alpha := 1.
	beta := 1.
	weights := map[*voronoi.Cell]float64{}
	zoneCentroid := utils.CellSetCentroid(m.zoneCells(zoneId))
	for c := range cs.Iter() {
		distanceFromZone := utils.Distance(c.Site, zoneCentroid)
		borderLength := 0.
		for _, he := range c.Halfedges {
			oc := he.Edge.GetOtherCell(c)
			if oc != nil && m.cellZone[oc] == zoneId {
				borderLength += utils.Distance(he.GetStartpoint(), he.GetEndpoint())
			}
		}
		if _, ok := weights[c]; !ok {
			weights[c] = alpha/distanceFromZone + beta*borderLength
		} else {
			weights[c] += alpha/distanceFromZone + beta*borderLength
		}
	}
	return weights
}

func addCells(cells []*voronoi.Cell, cs *voronoi.CellSet) {
	for _, c := range cells {
		cs.Add(c)
	}
}

func (m *Map) zoneError() float64 {
	e := 0.
	for z, cs := range m.zoneCell {
		da := utils.CellSetArea(cs) - m.c.ZoneSizes[z]
		e += da * da
	}
	return e
}

// func ZoneSizes(zc *ZoneConfig, zones map[*voronoi.Cell]int, d *voronoi.Diagram) (float64, []float64) {
// 	unassignedSize := 0.
// 	sizes := make([]float64, len(zc.Sizes))
// 	for _, c := range d.Cells {
// 		if z, ok := zones[c]; !ok {
// 			unassignedSize -= utils.CellArea(c)
// 		} else {
// 			sizes[z] -= utils.CellArea(c)
// 		}
// 	}
// 	return unassignedSize, sizes
// }

var palette = []color.Color{
	color.RGBA{0x9d, 0x9d, 0x9d, 0xff},
	color.RGBA{0xbe, 0x26, 0x33, 0xff},
	color.RGBA{0xe0, 0x6f, 0x8b, 0xff},
	color.RGBA{0x49, 0x3c, 0x2b, 0xff},
	color.RGBA{0xa4, 0x64, 0x22, 0xff},
	color.RGBA{0xeb, 0x89, 0x31, 0xff},
	color.RGBA{0xf7, 0xe2, 0x6b, 0xff},
	color.RGBA{0x2f, 0x48, 0x4e, 0xff},
	color.RGBA{0x44, 0x89, 0x1a, 0xff},
	color.RGBA{0xa3, 0xce, 0x27, 0xff},
	color.RGBA{0x1b, 0x26, 0x32, 0xff},
	color.RGBA{0x00, 0x57, 0x84, 0xff},
	color.RGBA{0x31, 0xa2, 0xf2, 0xff},
	color.RGBA{0xb2, 0xdc, 0xef, 0xff},
}

func (m *Map) saveImage(fileName string) {
	dest := image.NewRGBA(image.Rect(0, 0, 500, 500))
	gc := draw2dimg.NewGraphicContext(dest)
	scale := 500.
	xTrans := 0.
	yTrans := 0.

	// Zone fill
	for c, z := range m.cellZone {
		gc.SetFillColor(palette[z])
		gc.MoveTo(c.Halfedges[0].GetStartpoint().X*scale+xTrans, c.Halfedges[0].GetStartpoint().Y*scale+yTrans)
		for _, he := range c.Halfedges {
			gc.LineTo(he.GetEndpoint().X*scale+xTrans, he.GetEndpoint().Y*scale+yTrans)
		}
		gc.Close()
		gc.Fill()
	}
	gc.SetFillColor(image.Black)
	gc.SetFontSize(3)
	gc.SetLineWidth(0.1)

	// Zone centroids
	for c, _ := range m.cellZone {
		centroid := utils.CellCentroid(c)
		gc.MoveTo(centroid.X*scale+xTrans, centroid.Y*scale+yTrans)
		draw2dkit.Circle(gc, centroid.X*scale, centroid.Y*scale, 1)
	}

	// Zone edges
	gc.SetFillColor(image.Black)
	for _, e := range m.d.Edges {
		gc.MoveTo(e.Va.X*scale+xTrans, e.Va.Y*scale+yTrans)
		gc.LineTo(e.Vb.X*scale+xTrans, e.Vb.Y*scale+yTrans)
		gc.Close()
		gc.FillStroke()
	}

	// Save to file
	if err := draw2dimg.SaveToPngFile(fileName, dest); err != nil {
		log.Fatal(err)
	}
}
