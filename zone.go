package nrmg

import (
	"fmt"
	"image"
	"log"
	"math/rand"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dpdf"
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

func AssignZones(d *voronoi.Diagram, zc *ZoneConfig) map[*voronoi.Cell]int {
	numZones := len(zc.Sizes)
	zones := map[*voronoi.Cell]int{}
	zoneAreas := make([]float64, numZones)
	log.Printf("len(zoneAreas): %d", len(zoneAreas))
	eligibleCells := append([]*voronoi.Cell{}, d.Cells[rand.Intn(len(d.Cells))])
	currentZone := 0
	for {
		if len(eligibleCells) == 0 {
			break
		}
		log.Printf("current zone: %d", currentZone)
		ix := rand.Intn(len(eligibleCells))
		c := eligibleCells[ix]
		zones[c] = currentZone
		zoneAreas[currentZone] -= utils.CellArea(c)
		eligibleCells = append(eligibleCells[:ix], eligibleCells[ix+1:]...)
		for _, he := range c.Halfedges {
			oc := he.Edge.GetOtherCell(c)
			if oc != nil {
				if _, ok := zones[oc]; !ok {
					eligibleCells = append(eligibleCells, oc)
				}
			}
		}
		if zoneAreas[currentZone] > zc.Sizes[currentZone] {
			currentZone++
			eligibleCells = eligibleCells[:1]
		}
		if currentZone == numZones {
			break
		}
	}
	return zones
}

func RenderPreview(zones map[*voronoi.Cell]int, diagram *voronoi.Diagram) {
	// Initialize the graphic context on an RGBA image
	draw2d.SetFontFolder("resource/font")
	dest := draw2dpdf.NewPdf("L", "mm", "A4")
	gc := draw2dpdf.NewGraphicContext(dest)

	// Set the font luximbi.ttf
	// gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic})
	// Set the fill text color to black
	gc.SetFillColor(image.Black)
	gc.SetFontSize(3)
	gc.SetLineWidth(0.1)
	// Display Hello World
	scale := 150.
	for c, z := range zones {
		centroid := utils.CellCentroid(c)
		gc.MoveTo(centroid.X*scale, centroid.Y*scale)
		// draw2dkit.Circle(gc, centroid.X*scale, centroid.Y*scale, 1)
		s := fmt.Sprintf("%d", z)
		// l, t, r, b := gc.GetStringBounds(s)
		// log.Printf("string bounds %f, %f, %f, %f", l, t, r, b)
		gc.FillStringAt(s, centroid.X*scale-1, centroid.Y*scale-4.5)
	}

	for _, e := range diagram.Edges {
		gc.MoveTo(e.Va.X*scale, e.Va.Y*scale)
		gc.LineTo(e.Vb.X*scale, e.Vb.Y*scale)
		gc.Close()
		gc.FillStroke()
	}

	// Save to file
	if err := draw2dpdf.SaveToPdfFile("hello.pdf", dest); err != nil {
		log.Fatal(err)
	}
}

func Voronoi(n int) *voronoi.Diagram {
	bb := voronoi.NewBBox(0., 1., 0., 1.)
	return voronoi.ComputeDiagram(utils.RandomSites(bb, n), bb, true)
}
