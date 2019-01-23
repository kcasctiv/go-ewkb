package geo

type Point interface {
	X() float64
	Y() float64
	Z() float64
	M() float64
}

type point struct{ x, y float64 }

func (p point) X() float64 { return p.x }
func (p point) Y() float64 { return p.y }
func (p point) Z() float64 { return 0 }
func (p point) M() float64 { return 0 }

func NewPoint(x, y float64) Point {
	return point{x, y}
}

type pointZ struct{ x, y, z float64 }

func (p pointZ) X() float64 { return p.x }
func (p pointZ) Y() float64 { return p.y }
func (p pointZ) Z() float64 { return p.z }
func (p pointZ) M() float64 { return 0 }

func NewPointZ(x, y, z float64) Point {
	return pointZ{x, y, z}
}

type pointM struct{ x, y, m float64 }

func (p pointM) X() float64 { return p.x }
func (p pointM) Y() float64 { return p.y }
func (p pointM) Z() float64 { return 0 }
func (p pointM) M() float64 { return p.m }

func NewPointM(x, y, m float64) Point {
	return pointM{x, y, m}
}

type pointZM struct{ x, y, z, m float64 }

func (p pointZM) X() float64 { return p.x }
func (p pointZM) Y() float64 { return p.y }
func (p pointZM) Z() float64 { return p.z }
func (p pointZM) M() float64 { return p.m }

func NewPointZM(x, y, z, m float64) Point {
	return pointZM{x, y, z, m}
}

type MultiPoint interface {
	Point(int) Point
	Len() int
}

type multiPoint struct {
	points []Point
}

func (p multiPoint) Point(idx int) Point { return p.points[idx] }
func (p multiPoint) Len() int            { return len(p.points) }

func NewMultiPoint(points []Point) MultiPoint {
	return multiPoint{points}
}

type Polygon interface {
	Ring(int) MultiPoint
	Len() int
}

type polygon struct {
	rings []MultiPoint
}

func (p polygon) Ring(idx int) MultiPoint { return p.rings[idx] }
func (p polygon) Len() int                { return len(p.rings) }

func NewPolygon(rings []MultiPoint) Polygon {
	return &polygon{rings}
}

type MultiLine interface {
	Line(int) MultiPoint
	Len() int
}

type multiLine struct {
	lines []MultiPoint
}

func (l multiLine) Line(idx int) MultiPoint { return l.lines[idx] }
func (l multiLine) Len() int                { return len(l.lines) }

func NewMultiLine(lines []MultiPoint) MultiLine {
	return multiLine{lines}
}

type MultiPolygon interface {
	Polygon(int) Polygon
	Len() int
}

type multiPolygon struct {
	pols []Polygon
}

func (p multiPolygon) Polygon(idx int) Polygon { return p.pols[idx] }
func (p multiPolygon) Len() int                { return len(p.pols) }

func NewMultiPolygon(pols []Polygon) MultiPolygon {
	return multiPolygon{pols}
}
