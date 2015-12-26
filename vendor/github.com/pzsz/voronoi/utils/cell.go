// Copyright 2013 Przemyslaw Szczepaniak.
// MIT License: See https://github.com/gorhill/Javascript-Voronoi/LICENSE.md

// Author: Przemyslaw Szczepaniak (przeszczep@gmail.com)
// Utils for processing voronoi diagrams

package utils

import "github.com/pzsz/voronoi"

// Calculate area of a cell
func CellArea(cell *voronoi.Cell) float64 {
	area := 0.
	for _, halfedge := range cell.Halfedges {
		s := halfedge.GetStartpoint()
		e := halfedge.GetEndpoint()
		area += s.X * e.Y
		area -= s.Y * e.X
	}
	return area / 2
}

func CellSetArea(cs *voronoi.CellSet) float64 {
	area := 0.
	for c := range cs.Iter() {
		area -= CellArea(c)
	}
	return area
}

func AreNeighbors(first, second *voronoi.Cell) bool {
	for _, he := range first.Halfedges {
		if he.Edge.GetOtherCell(first) == second {
			return true
		}
	}
	return false
}

// Calculate centroid of a cell
func CellCentroid(cell *voronoi.Cell) voronoi.Vertex {
	x, y := float64(0), float64(0)
	for _, halfedge := range cell.Halfedges {
		s := halfedge.GetStartpoint()
		e := halfedge.GetEndpoint()
		v := s.X*e.Y - e.X*s.Y
		x += (s.X + e.X) * v
		y += (s.Y + e.Y) * v
	}
	v := CellArea(cell) * 6
	return voronoi.Vertex{X: x / v, Y: y / v}
}

func CellSetCentroid(cs *voronoi.CellSet) voronoi.Vertex {
	x, y := float64(0), float64(0)
	for c := range cs.Iter() {
		cc := CellCentroid(c)
		x += cc.X
		y += cc.Y
	}
	v := float64(cs.Cardinality())
	return voronoi.Vertex{X: x / v, Y: y / v}
}

// Calculate centroid of a cell
func InsideCell(cell *voronoi.Cell, v voronoi.Vertex) bool {
	for _, halfedge := range cell.Halfedges {
		a := halfedge.GetStartpoint()
		b := halfedge.GetEndpoint()

		cross := ((b.X-a.X)*(v.Y-a.Y) - (b.Y-a.Y)*(v.X-a.X))

		if cross > 0 {
			return false
		}
	}
	return true
}

func EdgeIndex(cell *voronoi.Cell, edge *voronoi.Edge) int {
	for i, halfedge := range cell.Halfedges {
		if halfedge.Edge == edge {
			return i
		}
	}
	return -1
}
