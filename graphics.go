package nrmg

import (
	"image"
	"image/color"
	"log"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/pzsz/voronoi/utils"
)

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
