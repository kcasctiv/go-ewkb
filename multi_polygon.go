package ewkb

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/kcasctiv/go-ewkb/geo"
)

// MultiPolygon presents MultiPolygon geometry object
type MultiPolygon struct {
	header
	mp geo.MultiPolygon
}

// Polygon returns polygon with specified index
func (p *MultiPolygon) Polygon(idx int) geo.Polygon { return p.mp.Polygon(idx) }

// Len returns count of polygons
func (p *MultiPolygon) Len() int { return p.mp.Len() }

// String returns WKT/EWKT geometry representation
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

// Scan implements sql.Scanner interface
func (p *MultiPolygon) Scan(src interface{}) error {
	return scanGeometry(src, p)
}

// Value implements sql driver.Valuer interface
func (p *MultiPolygon) Value() (driver.Value, error) {
	return p.String(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface
func (p *MultiPolygon) UnmarshalBinary(data []byte) error {
	h, byteOrder, offset := readHeader(data)
	if h.Type() != MultiPolygonType {
		return errors.New("not expected geometry type")
	}

	p.header = h

	var err error
	p.mp, _, err = readMultiPolygon(
		data[offset:],
		byteOrder,
		getReadPointFunc(h.wkbType),
	)
	return err
}
