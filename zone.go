package nrmg

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/persomi/set"
	"github.com/pzsz/voronoi"
	"github.com/pzsz/voronoi/utils"
)

type Zone struct {
	id int
}

type ZoneConfig struct {
	Sizes []float64
}

func (zc *ZoneConfig) NormalizeSizes() {
	total := 0.
	for _, s := range zc.Sizes {
		total += s
	}
	for ix, s := range zc.Sizes {
		zc.Sizes[ix] = s / total
	}
}

func (zc ZoneConfig) String() string {
	var out string
	for ix, s := range zc.Sizes {
		out += fmt.Sprintf(" zone %d: %f |", ix, s)
	}
	return out
}

var Palette = []color.Color{
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

func ZoneSizes(zc *ZoneConfig, zones map[*voronoi.Cell]int, d *voronoi.Diagram) (float64, []float64) {
	unassignedSize := 0.
	sizes := make([]float64, len(zc.Sizes))
	for _, c := range d.Cells {
		if z, ok := zones[c]; !ok {
			unassignedSize -= utils.CellArea(c)
		} else {
			sizes[z] -= utils.CellArea(c)
		}
	}
	return unassignedSize, sizes
}

type Map struct {
	d         *voronoi.Diagram
	cellZone  map[*voronoi.Cell]int
	zoneCell  map[int]cellSet
	zc        *ZoneConfig
	zoneCount int
}

func New() *Map {
	return &Map{
		cellZone: map[*voronoi.Cell]int{},
		zoneCell: map[int]cellSet{},
	}
}

type cellSet set.Interface

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
	out := set.New(set.NonThreadSafe)
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
	for i := 0; i <= m.zoneCount; i++ {
		if i != zoneId {
			out = append(out, i)
		}
	}
	return out
}

func randomChoice(cs cellSet) *voronoi.Cell {
	list := cs.List()
	ix := rand.Intn(len(list))
	return list[ix].(*voronoi.Cell)
}

func (m *Map) growZone(zoneId int, startCell *voronoi.Cell) (cellSet, error) {
	var copiedMap Map = *m
	area := 0.
	prescribedArea := m.zc.Sizes[zoneId]
	grow := set.New(set.NonThreadSafe)
	grow.AddOne(startCell)
	for {
		if area >= prescribedArea {
			m.zoneCell = copiedMap.zoneCell
			m.cellZone = copiedMap.cellZone
			return grow, nil
		}
		if grow.IsEmpty() {
			return copiedMap.zoneCells(zoneId), fmt.Errorf("No space for zone %d", zoneId)
		}
		next := randomChoice(grow)
		grow.Remove(next)
		// log.Println(grow)
		copiedMap.assign(zoneId, next)
		area -= utils.CellArea((*voronoi.Cell)(next))
		grow.Merge(copiedMap.neighbors(next))
	}
}

func addCells(cells []*voronoi.Cell, cs cellSet) {
	for _, c := range cells {
		cs.AddOne(c)
	}
}

func (m *Map) AssignZones(d *voronoi.Diagram, zc *ZoneConfig) {
	m.d = d
	m.zc = zc
	m.zoneCount = len(zc.Sizes)
	grow := set.New(set.NonThreadSafe)
	addCells(m.d.Cells, grow)
	for zoneId := 0; zoneId < m.zoneCount; zoneId++ {
		start := randomChoice(grow)
		cs, err := m.growZone(zoneId, start)
		if err != nil {
			growOther := m.zoneNeighbors(m.otherZones(zoneId)...)
			growOther.Separate(cs)
			start := randomChoice(grow)
			cs, err := m.growZone(zoneId, start)
			if err != nil {
				log.Printf("exiting assignment: %v", err)
				return
			}
			grow = cs
			continue
		}
		grow = cs
	}
}

func (m *Map) SaveImage(fileName string) {
	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, 500, 500))
	gc := draw2dimg.NewGraphicContext(dest)
	// Set the font luximbi.ttf
	// gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic})
	// Set the fill text color to black

	scale := 500.
	xTrans := 0.
	yTrans := 0.

	for c, z := range m.cellZone {
		gc.SetFillColor(Palette[z+5])
		gc.MoveTo(c.Halfedges[0].GetStartpoint().X*scale+xTrans, c.Halfedges[0].GetStartpoint().Y*scale+yTrans)
		for _, he := range c.Halfedges {
			// gc.MoveTo(he.GetStartpoint().X*scale, he.GetStartpoint().Y*scale)
			gc.LineTo(he.GetEndpoint().X*scale+xTrans, he.GetEndpoint().Y*scale+yTrans)
		}
		gc.Close()
		gc.Fill()
	}

	gc.SetFillColor(image.Black)
	gc.SetFontSize(3)
	gc.SetLineWidth(0.1)

	// for c, z := range zones {
	// 	centroid := utils.CellCentroid(c)
	// 	gc.MoveTo(centroid.X*scale+xTrans, centroid.Y*scale+yTrans)
	// 	// draw2dkit.Circle(gc, centroid.X*scale, centroid.Y*scale, 1)
	// 	s := fmt.Sprintf("%d", z)
	// 	// l, t, r, b := gc.GetStringBounds(s)
	// 	// log.Printf("string bounds %f, %f, %f, %f", l, t, r, b)
	// 	gc.FillStringAt(s, centroid.X*scale-1+xTrans, centroid.Y*scale-4.5+yTrans)
	// }

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

func Voronoi(n int) *voronoi.Diagram {
	bb := voronoi.NewBBox(0., 1., 0., 1.)
	sites := utils.RandomSites(bb, n)
	d := voronoi.ComputeDiagram(sites, bb, true)
	for i := 0; i < 1; i++ {
		sites = utils.LloydRelaxation(d.Cells)
		d = voronoi.ComputeDiagram(sites, bb, true)
	}
	return d
}
