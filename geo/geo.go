// Package geo contains interfaces
// of simple geometry objects,
// and their creation methods
package geo

// Point presents interface of point
type Point interface {
	// X returns value of X dimension
	X() float64
	// Y returns value of X dimension
	Y() float64
	// Z returns value of X dimension
	Z() float64
	// M returns value of X dimension
	M() float64
}

type point struct{ x, y float64 }

func (p point) X() float64 { return p.x }
func (p point) Y() float64 { return p.y }
func (p point) Z() float64 { return 0 }
func (p point) M() float64 { return 0 }

// NewPoint returns new 2 dimensions point
func NewPoint(x, y float64) Point {
	return point{x, y}
}

type pointZ struct{ x, y, z float64 }

func (p pointZ) X() float64 { return p.x }
func (p pointZ) Y() float64 { return p.y }
func (p pointZ) Z() float64 { return p.z }
func (p pointZ) M() float64 { return 0 }

// NewPointZ returns new 3 dimensions point with Z dimension
func NewPointZ(x, y, z float64) Point {
	return pointZ{x, y, z}
}

type pointM struct{ x, y, m float64 }

func (p pointM) X() float64 { return p.x }
func (p pointM) Y() float64 { return p.y }
func (p pointM) Z() float64 { return 0 }
func (p pointM) M() float64 { return p.m }

// NewPointM returns new 3 dimensions point with M dimension
func NewPointM(x, y, m float64) Point {
	return pointM{x, y, m}
}

type pointZM struct{ x, y, z, m float64 }

func (p pointZM) X() float64 { return p.x }
func (p pointZM) Y() float64 { return p.y }
func (p pointZM) Z() float64 { return p.z }
func (p pointZM) M() float64 { return p.m }

// NewPointZM returns new 4 dimensions point
func NewPointZM(x, y, z, m float64) Point {
	return pointZM{x, y, z, m}
}

// MultiPoint presents interface of multi point
type MultiPoint interface {
	// Point returns point with specified index
	Point(int) Point
	// Len returns count of points
	Len() int
}

type multiPoint struct {
	points []Point
}

func (p multiPoint) Point(idx int) Point { return p.points[idx] }
func (p multiPoint) Len() int            { return len(p.points) }

// NewMultiPoint returns new multi point
func NewMultiPoint(points []Point) MultiPoint {
	return multiPoint{points}
}

// Polygon presents interface of polygon
type Polygon interface {
	// Ring returns ring with specified index
	Ring(int) MultiPoint
	// Len returns count of rings
	Len() int
}

type polygon struct {
	rings []MultiPoint
}

func (p polygon) Ring(idx int) MultiPoint { return p.rings[idx] }
func (p polygon) Len() int                { return len(p.rings) }

// NewPolygon returns new polygon
func NewPolygon(rings []MultiPoint) Polygon {
	return polygon{rings}
}

// MultiLine presents interface of multi line
type MultiLine interface {
	// Line returns line with specified index
	Line(int) MultiPoint
	// Len returns count of lines
	Len() int
}

type multiLine struct {
	lines []MultiPoint
}

func (l multiLine) Line(idx int) MultiPoint { return l.lines[idx] }
func (l multiLine) Len() int                { return len(l.lines) }

// NewMultiLine returns new multi line
func NewMultiLine(lines []MultiPoint) MultiLine {
	return multiLine{lines}
}

// MultiPolygon presents interface of multi polygon
type MultiPolygon interface {
	// Polygon returns polygon with specified index
	Polygon(int) Polygon
	// Len returns count of polygons
	Len() int
}

type multiPolygon struct {
	pols []Polygon
}

func (p multiPolygon) Polygon(idx int) Polygon { return p.pols[idx] }
func (p multiPolygon) Len() int                { return len(p.pols) }

// NewMultiPolygon returns new multi polygon
func NewMultiPolygon(pols []Polygon) MultiPolygon {
	return multiPolygon{pols}
}
