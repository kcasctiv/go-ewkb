package ewkb

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/kcasctiv/go-ewkb/geo"
)

// Polygon presents Polygon geometry object
type Polygon struct {
	header
	poly geo.Polygon
}

// NewPolygon returns new Polygon,
// created from geometry base and coords data
func NewPolygon(b Base, poly geo.Polygon) Polygon {
	return Polygon{
		header: header{
			byteOrder: b.ByteOrder(),
			wkbType: getFlags(
				b.HasZ(),
				b.HasM(),
				b.HasSRID(),
			) | PolygonType,
			srid: b.SRID(),
		},
		poly: poly,
	}
}

// Ring returns ring with specified index
func (p *Polygon) Ring(idx int) geo.MultiPoint { return p.poly.Ring(idx) }

// Len returns count of rings
func (p *Polygon) Len() int { return p.poly.Len() }

// String returns WKT/EWKT geometry representation
func (p *Polygon) String() string {
	var s string
	if p.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", p.srid)
	}
	s += "POLYGON"
	if !p.HasZ() && p.HasM() {
		s += "M"
	}

	return s + printPolygon(p, p.HasZ(), p.HasM())
}

// Scan implements sql.Scanner interface
func (p *Polygon) Scan(src interface{}) error {
	return scanGeometry(src, p)
}

// Value implements sql driver.Valuer interface
func (p *Polygon) Value() (driver.Value, error) {
	return p.String(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface
func (p *Polygon) UnmarshalBinary(data []byte) error {
	h, byteOrder, offset := readHeader(data)
	if h.Type() != PolygonType {
		return errors.New("not expected geometry type")
	}

	p.header = h

	var err error
	p.poly, _, err = readPolygon(data[offset:], byteOrder, getReadPointFunc(h.wkbType))
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler interface
func (p *Polygon) MarshalBinary() ([]byte, error) {
	size := headerSize(p.HasSRID()) + polygonSize(p, p.HasZ(), p.HasM())
	b := make([]byte, size)

	byteOrder := getBinaryByteOrder(p.ByteOrder())
	offset := writeHeader(p, p.Type(), byteOrder, p.HasSRID(), b)
	writePolygon(p, byteOrder, p.HasZ(), p.HasM(), b[offset:])

	return b, nil
}

func printPolygon(p geo.Polygon, hasZ, hasM bool) string {
	if p.Len() == 0 {
		return " EMPTY"
	}

	var s string
	for idx := 0; idx < p.Len(); idx++ {
		s += printMultiPoint(p.Ring(idx), hasZ, hasM) + ","
	}

	return "(" + s[:len(s)-1] + ")"
}
