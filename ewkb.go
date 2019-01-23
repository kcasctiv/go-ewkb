package ewkb

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
	"math"

	"github.com/kcasctiv/go-ewkb/geo"
)

// Byte orders
const (
	XDR byte = 0x00 // Big endian
	NDR byte = 0x01 // Little endian
)

// Type flags
const (
	zFlag    uint32 = 0x80000000 // Z dimension flag
	mFlag    uint32 = 0x40000000 // M dimension flag
	sridFlag uint32 = 0x20000000 // SRID flag
	bboxFlag uint32 = 0x10000000 // BBOX flag
)

// Point presents 2, 3 or 4 dimensions point
type Point struct {
	byteOrder  byte
	wkbType    uint32
	srid       int32
	bbox       *bbox
	x, y, z, m float64
}

func (p *Point) ByteOrder() byte { return p.byteOrder }
func (p *Point) Type() uint32    { return p.wkbType & uint32(math.MaxUint16) }
func (p *Point) HasZ() bool      { return (p.wkbType & zFlag) == zFlag }
func (p *Point) HasM() bool      { return (p.wkbType & mFlag) == mFlag }
func (p *Point) HasSRID() bool   { return (p.wkbType & sridFlag) == sridFlag }
func (p *Point) HasBBOX() bool   { return (p.wkbType & bboxFlag) == bboxFlag }
func (p *Point) X() float64      { return p.x }
func (p *Point) Y() float64      { return p.y }
func (p *Point) Z() float64      { return p.z }
func (p *Point) M() float64      { return p.m }

func (p *Point) String() string {
	var s string
	if p.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", p.srid)
	}
	s += "POINT "
	if p.HasZ() {
		s += "Z"
	}
	if p.HasM() {
		s += "M"
	}

	return s + " (" + printPoint(p, p.HasZ(), p.HasM()) + ")"
}

func (p *Point) Scan(src interface{}) error {
	// TODO:
	return nil
}

func (p *Point) Value() (driver.Value, error) {
	return p.String(), nil
}

func printPoint(p geo.Point, hasZ, hasM bool) string {
	s := fmt.Sprintf("%f %f", p.X(), p.Y())
	if hasZ {
		s += fmt.Sprintf(" %f", p.Z())
	}
	if hasM {
		s += fmt.Sprintf(" %f", p.M())
	}

	return s
}

type LineString struct {
	byteOrder byte
	wkbType   uint32
	srid      int32
	bbox      *bbox
	mp        geo.MultiPoint
}

func (l *LineString) ByteOrder() byte         { return l.byteOrder }
func (l *LineString) Type() uint32            { return l.wkbType & uint32(math.MaxUint16) }
func (l *LineString) HasZ() bool              { return (l.wkbType & zFlag) == zFlag }
func (l *LineString) HasM() bool              { return (l.wkbType & mFlag) == mFlag }
func (l *LineString) HasSRID() bool           { return (l.wkbType & sridFlag) == sridFlag }
func (l *LineString) HasBBOX() bool           { return (l.wkbType & bboxFlag) == bboxFlag }
func (l *LineString) Point(idx int) geo.Point { return l.mp.Point(idx) }
func (l *LineString) Len() int                { return l.mp.Len() }

func (l *LineString) String() string {
	var s string
	if l.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", l.srid)
	}
	s += "LINESTRING "
	if l.HasZ() {
		s += "Z"
	}
	if l.HasM() {
		s += "M"
	}

	return s + " " + printMultiPoint(l, l.HasZ(), l.HasM())
}

func (l *LineString) Scan(src interface{}) error {
	// TODO:
	return nil
}

func (l *LineString) Value() (driver.Value, error) {
	return l.String(), nil
}

func printMultiPoint(p geo.MultiPoint, hasZ, hasM bool) string {
	if p.Len() == 0 {
		return "()"
	}

	var s string
	for idx := 0; idx < p.Len(); idx++ {
		s += printPoint(p.Point(idx), hasZ, hasM) + ", "
	}

	return "(" + s[:len(s)-2] + ")"
}

type Polygon struct {
	byteOrder byte
	wkbType   uint32
	srid      int32
	bbox      *bbox
	poly      geo.Polygon
}

func (p *Polygon) ByteOrder() byte             { return p.byteOrder }
func (p *Polygon) Type() uint32                { return p.wkbType & uint32(math.MaxUint16) }
func (p *Polygon) HasZ() bool                  { return (p.wkbType & zFlag) == zFlag }
func (p *Polygon) HasM() bool                  { return (p.wkbType & mFlag) == mFlag }
func (p *Polygon) HasSRID() bool               { return (p.wkbType & sridFlag) == sridFlag }
func (p *Polygon) HasBBOX() bool               { return (p.wkbType & bboxFlag) == bboxFlag }
func (p *Polygon) Ring(idx int) geo.MultiPoint { return p.poly.Ring(idx) }
func (p *Polygon) Len() int                    { return p.poly.Len() }

func (p *Polygon) String() string {
	var s string
	if p.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", p.srid)
	}
	s += "POLYGON "
	if p.HasZ() {
		s += "Z"
	}
	if p.HasM() {
		s += "M"
	}

	return s + " " + printPolygon(p, p.HasZ(), p.HasM())
}

func (p *Polygon) Scan(src interface{}) error {
	// TODO:
	return nil
}

func (p *Polygon) Value() (driver.Value, error) {
	return p.String(), nil
}

func printPolygon(p geo.Polygon, hasZ, hasM bool) string {
	if p.Len() == 0 {
		return "()"
	}

	var s string
	for idx := 0; idx < p.Len(); idx++ {
		s += printMultiPoint(p.Ring(idx), hasZ, hasM) + ", "
	}

	return "(" + s[:len(s)-2] + ")"
}

type MultiPoint struct {
	byteOrder byte
	wkbType   uint32
	srid      int32
	bbox      *bbox
	mp        geo.MultiPoint
}

func (p *MultiPoint) ByteOrder() byte         { return p.byteOrder }
func (p *MultiPoint) Type() uint32            { return p.wkbType & uint32(math.MaxUint16) }
func (p *MultiPoint) HasZ() bool              { return (p.wkbType & zFlag) == zFlag }
func (p *MultiPoint) HasM() bool              { return (p.wkbType & mFlag) == mFlag }
func (p *MultiPoint) HasSRID() bool           { return (p.wkbType & sridFlag) == sridFlag }
func (p *MultiPoint) HasBBOX() bool           { return (p.wkbType & bboxFlag) == bboxFlag }
func (p *MultiPoint) Point(idx int) geo.Point { return p.mp.Point(idx) }
func (p *MultiPoint) Len() int                { return p.mp.Len() }

func (p *MultiPoint) String() string {
	var s string
	if p.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", p.srid)
	}
	s += "MULTIPOINT "
	if p.HasZ() {
		s += "Z"
	}
	if p.HasM() {
		s += "M"
	}

	return s + " " + printMultiPoint(p, p.HasZ(), p.HasM())
}

func (p *MultiPoint) Scan(src interface{}) error {
	// TODO:
	return nil
}

func (p *MultiPoint) Value() (driver.Value, error) {
	return p.String(), nil
}

type MultiLineString struct {
	byteOrder byte
	wkbType   uint32
	srid      int32
	bbox      *bbox
	ml        geo.MultiLine
}

func (l *MultiLineString) ByteOrder() byte             { return l.byteOrder }
func (l *MultiLineString) Type() uint32                { return l.wkbType & uint32(math.MaxUint16) }
func (l *MultiLineString) HasZ() bool                  { return (l.wkbType & zFlag) == zFlag }
func (l *MultiLineString) HasM() bool                  { return (l.wkbType & mFlag) == mFlag }
func (l *MultiLineString) HasSRID() bool               { return (l.wkbType & sridFlag) == sridFlag }
func (l *MultiLineString) HasBBOX() bool               { return (l.wkbType & bboxFlag) == bboxFlag }
func (l *MultiLineString) Line(idx int) geo.MultiPoint { return l.ml.Line(idx) }
func (l *MultiLineString) Len() int                    { return l.ml.Len() }

func (l *MultiLineString) String() string {
	var s string
	if l.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", l.srid)
	}
	s += "MULTILINESTRING "
	if l.HasZ() {
		s += "Z"
	}
	if l.HasM() {
		s += "M"
	}
	s += " ("
	if l.Len() > 0 {
		for idx := 0; idx < l.Len(); idx++ {
			s += printMultiPoint(l.Line(idx), l.HasZ(), l.HasM()) + ", "
		}

		s = s[:len(s)-2]
	}

	return s + ")"
}

func (l *MultiLineString) Scan(src interface{}) error {
	// TODO:
	return nil
}

func (l *MultiLineString) Value() (driver.Value, error) {
	return l.String(), nil
}

type MultiPolygon struct {
	byteOrder byte
	wkbType   uint32
	srid      int32
	bbox      *bbox
	mp        geo.MultiPolygon
}

func (p *MultiPolygon) ByteOrder() byte             { return p.byteOrder }
func (p *MultiPolygon) Type() uint32                { return p.wkbType & uint32(math.MaxUint16) }
func (p *MultiPolygon) HasZ() bool                  { return (p.wkbType & zFlag) == zFlag }
func (p *MultiPolygon) HasM() bool                  { return (p.wkbType & mFlag) == mFlag }
func (p *MultiPolygon) HasSRID() bool               { return (p.wkbType & sridFlag) == sridFlag }
func (p *MultiPolygon) HasBBOX() bool               { return (p.wkbType & bboxFlag) == bboxFlag }
func (p *MultiPolygon) Polygon(idx int) geo.Polygon { return p.mp.Polygon(idx) }
func (p *MultiPolygon) Len() int                    { return p.mp.Len() }

func (p *MultiPolygon) String() string {
	var s string
	if p.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", p.srid)
	}
	s += "MULTIPOLYGON "
	if p.HasZ() {
		s += "Z"
	}
	if p.HasM() {
		s += "M"
	}
	s += " ("
	if p.Len() > 0 {
		for idx := 0; idx < p.Len(); idx++ {
			s += printPolygon(p.Polygon(idx), p.HasZ(), p.HasM()) + ", "
		}

		s = s[:len(s)-2]
	}

	return s + ")"
}

func (p *MultiPolygon) Scan(src interface{}) error {
	// TODO:
	return nil
}

func (p *MultiPolygon) Value() (driver.Value, error) {
	return p.String(), nil
}

type Wrapper struct {
	geom Geometry
}

func (w *Wrapper) ByteOrder() byte              { return w.geom.ByteOrder() }
func (w *Wrapper) Type() uint32                 { return w.geom.Type() }
func (w *Wrapper) HasZ() bool                   { return w.geom.HasZ() }
func (w *Wrapper) HasM() bool                   { return w.geom.HasM() }
func (w *Wrapper) HasSRID() bool                { return w.geom.HasSRID() }
func (w *Wrapper) HasBBOX() bool                { return w.geom.HasBBOX() }
func (w *Wrapper) String() string               { return w.geom.String() }
func (w *Wrapper) Value() (driver.Value, error) { return w.geom.Value() }
func (w *Wrapper) Geometry() Geometry           { return w.geom }

func (w *Wrapper) Scan(src interface{}) error {
	// TODO:
	return nil
}

type Geometry interface {
	ByteOrder() byte
	Type() uint32
	HasZ() bool
	HasM() bool
	HasSRID() bool
	HasBBOX() bool
	fmt.Stringer
	sql.Scanner
	driver.Valuer
	// TODO: implement these interfaces
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type bbox struct {
	xmin, xmax float64
	ymin, ymax float64
	zmin, zmax float64
	mmin, mmax float64
}
