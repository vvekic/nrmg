package nrmg

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/pzsz/voronoi"
	"github.com/pzsz/voronoi/utils"
)

type Zone struct {
	id int
}

type ZoneConfig struct {
	Sizes []float64
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

func (zc *ZoneConfig) NormalizeSizes() {
	total := 0.
	for _, s := range zc.Sizes {
		total += s
	}
	for ix, s := range zc.Sizes {
		zc.Sizes[ix] = s / total
	}
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

func AssignZones(d *voronoi.Diagram, zc *ZoneConfig) map[*voronoi.Cell]int {
	numZones := len(zc.Sizes)
	zones := map[*voronoi.Cell]int{}
	zoneAreas := make([]float64, numZones)
	// log.Printf("len(zoneAreas): %d", len(zoneAreas))
	eligibleCells := append([]*voronoi.Cell{}, d.Cells[rand.Intn(len(d.Cells))])
	var safeEligibleCells []*voronoi.Cell
	var safeZones map[*voronoi.Cell]int
	var safeZoneAreas []float64
	currentZone := 0
	i := 0
	maxCurrentSize := 0.
	var lastGoodZones map[*voronoi.Cell]int
	for {
		// RenderPreview(zones, d, fmt.Sprintf("assign/%d.png", i))
		i++
		if len(eligibleCells) == 0 {
			if len(safeEligibleCells) == 1 {
				return lastGoodZones
			}
			if zoneAreas[currentZone] > maxCurrentSize {
				maxCurrentSize = zoneAreas[currentZone]
				lastGoodZones = map[*voronoi.Cell]int{}
				for c, z := range zones {
					lastGoodZones[c] = z
				}
			}
			if len(safeEligibleCells) < 1 {
				break
			}
			safeEligibleCells = append([]*voronoi.Cell{}, safeEligibleCells[1:]...)
			eligibleCells = append([]*voronoi.Cell{}, safeEligibleCells[:1]...)
			zones = map[*voronoi.Cell]int{}
			for c, z := range safeZones {
				zones[c] = z
			}
			zoneAreas = append([]float64{}, safeZoneAreas...)
		}
		// log.Printf("current zone: %d", currentZone)
		ix := rand.Intn(len(eligibleCells))
		c := eligibleCells[ix]
		zones[c] = currentZone
		zoneAreas[currentZone] -= utils.CellArea(c)
		eligibleCells = append(eligibleCells[:ix], eligibleCells[ix+1:]...)
	loop:
		for _, he := range c.Halfedges {
			if utils.Distance(he.GetStartpoint(), he.GetEndpoint()) < 0.03 {
				continue loop
			}
			oc := he.Edge.GetOtherCell(c)
			if oc != nil {
				if _, ok := zones[oc]; !ok {
					for _, ec := range eligibleCells {
						if ec == oc {
							continue loop
						}
					}
					eligibleCells = append(eligibleCells, oc)
				}
			}
		}
		if zoneAreas[currentZone] > zc.Sizes[currentZone] {
			currentZone++
			maxCurrentSize = 0
			if len(eligibleCells) > 1 {
				safeEligibleCells = append([]*voronoi.Cell{}, eligibleCells...)
				// append other neighbor cells of other zones
				otherEligibleCells := []*voronoi.Cell{}
				for c, z := range zones {
					if z < currentZone {
					loopy:
						for _, he := range c.Halfedges {
							oc := he.Edge.GetOtherCell(c)
							if oc != nil {
								if _, ok := zones[oc]; !ok {
									for _, ec := range safeEligibleCells {
										if ec == oc {
											continue loopy
										}
									}
									for _, ec := range otherEligibleCells {
										if ec == oc {
											continue loopy
										}
									}
									otherEligibleCells = append(otherEligibleCells, oc)
								}
							}
						}
					}
				}
				//
				safeEligibleCells = append(safeEligibleCells, otherEligibleCells...)
				safeZones = map[*voronoi.Cell]int{}
				for c, z := range zones {
					safeZones[c] = z
				}
				safeZoneAreas = append([]float64{}, zoneAreas...)
			}
			eligibleCells = eligibleCells[:1]
		}
		if currentZone == numZones {
			break
		}
	}
	return zones
}

func RenderPreview(zones map[*voronoi.Cell]int, diagram *voronoi.Diagram, fileName string) {
	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, 500, 500))
	gc := draw2dimg.NewGraphicContext(dest)
	// Set the font luximbi.ttf
	// gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic})
	// Set the fill text color to black

	scale := 500.
	xTrans := 0.
	yTrans := 0.

	for _, c := range diagram.Cells {
		if _, ok := zones[c]; !ok {
			continue
		}
		gc.SetFillColor(Palette[zones[c]+5])
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
	for _, e := range diagram.Edges {
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
