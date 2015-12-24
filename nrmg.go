package nrmg

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/persomi/set"
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
	zoneCell  map[int]cellSet
	c         *Config
	zoneCount int
}

func New() *Map {
	return &Map{
		cellZone: map[*voronoi.Cell]int{},
		zoneCell: map[int]cellSet{},
	}
}

type cellSet set.Interface

func newCellSet() cellSet {
	return set.New(set.NonThreadSafe)
}

func (m *Map) tesselate(c *Config) {
	m.d = createDiagram(c.Cells)
	m.c = c
	m.zoneCount = len(c.ZoneSizes)
	grow := set.New(set.NonThreadSafe)
	addCells(m.d.Cells, grow)
	for zoneId := 0; zoneId < m.zoneCount; zoneId++ {
		grown, err := m.tryGrowingZone(zoneId, grow)
		if err != nil {
			grow = m.zoneNeighbors(m.otherZones(zoneId)...)
			grow.(set.Interface).Separate(grown)
			grown, err = m.tryGrowingZone(zoneId, grow)
			if err != nil {
				log.Printf("Zone %d error: %v", zoneId, err)
				continue
			}
		}
		m.assignSet(zoneId, grown)
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

func (m *Map) tryGrowingZone(zoneId int, grow cellSet) (cellSet, error) {
	for {
		origin := randomChoice(grow)
		grown, err := m.growZone(zoneId, origin)
		if err != nil {
			grow.(set.Interface).Separate(grown)
			if grow.(set.Interface).IsEmpty() {
				return grown, err
			}
			continue
		}
		return grown, nil
	}
}

func (m *Map) growZone(zoneId int, startCell *voronoi.Cell) (cellSet, error) {
	var copiedMap Map = *m
	area := 0.
	prescribedArea := m.c.ZoneSizes[zoneId]
	grow := set.New(set.NonThreadSafe)
	grow.AddOne(startCell)
	grown := set.New(set.NonThreadSafe)
	for {
		if area > prescribedArea {
			return grown, nil
		}
		if grow.IsEmpty() {
			return grown, fmt.Errorf("No space for zone %d", zoneId)
		}
		next := randomChoice(grow)
		grow.Remove(next)
		// log.Println(grow)
		copiedMap.assign(zoneId, next)
		grown.AddOne(next)
		area -= utils.CellArea((*voronoi.Cell)(next))
		grow.Merge(copiedMap.neighbors(next))
	}
}

func (m *Map) assignSet(zoneId int, csd cellSet) {
	for _, c := range csd.(set.Interface).List() {
		m.assign(zoneId, c.(*voronoi.Cell))
	}
}

func (m *Map) assign(zoneId int, cells ...*voronoi.Cell) {
	for _, c := range cells {
		if _, ok := m.cellZone[c]; !ok {
			m.cellZone[c] = zoneId
			if m.zoneCell[zoneId] == nil {
				m.zoneCell[zoneId] = set.New(set.NonThreadSafe)
			}
			m.zoneCell[zoneId].Add(c)
		} else {
			log.Fatalf("Cell already assigned: site: %v, zone: %d", c.Site, zoneId)
		}
	}
}

func (m *Map) zoneCells(zoneId int) cellSet {
	if c, ok := m.zoneCell[zoneId]; !ok {
		return set.New(set.NonThreadSafe)
	} else {
		return c
	}
}

func (m *Map) zoneNeighbors(zoneIds ...int) cellSet {
	out := set.New(set.NonThreadSafe)
	for _, zoneId := range zoneIds {

		zoneCells := m.zoneCells(zoneId).List()
		for _, c := range zoneCells {
			out.Merge(m.neighbors(c.(*voronoi.Cell)))
		}
	}
	return out
}

func (m *Map) neighbors(cells ...*voronoi.Cell) cellSet {
	out := newCellSet()
	for _, c := range cells {
		ocSet := set.New(set.NonThreadSafe)
		for _, he := range c.Halfedges {
			oc := he.Edge.GetOtherCell((*voronoi.Cell)(c))
			if oc != nil {
				if _, ok := m.cellZone[(*voronoi.Cell)(oc)]; !ok {
					ocSet.AddOne(oc)
				}
			}
		}
		out.Merge(ocSet)
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

func randomChoice(cs cellSet) *voronoi.Cell {
	list := cs.List()
	return list[rand.Intn(len(list))].(*voronoi.Cell)
}

func addCells(cells []*voronoi.Cell, cs cellSet) {
	for _, c := range cells {
		cs.AddOne(c)
	}
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
	color.RGBA{0xff, 0xff, 0xff, 0xff},
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
		gc.SetFillColor(palette[z+5])
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
