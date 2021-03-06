package ewkb

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/kcasctiv/go-ewkb/geo"
)

// MultiPoint presents MultiPoint geometry object
type MultiPoint struct {
	header
	mp geo.MultiPoint
}

// NewMultiPoint returns new MultiPoint,
// created from geometry base and coords data
func NewMultiPoint(b Base, mp geo.MultiPoint) MultiPoint {
	return MultiPoint{
		header: header{
			byteOrder: b.ByteOrder(),
			wkbType: getFlags(
				b.HasZ(),
				b.HasM(),
				b.HasSRID(),
			) | MultiPointType,
			srid: b.SRID(),
		},
		mp: mp,
	}
}

// Point returns point with specified index
func (p *MultiPoint) Point(idx int) geo.Point { return p.mp.Point(idx) }

// Len returns length of MultiPoint (count of points)
func (p *MultiPoint) Len() int { return p.mp.Len() }

// String returns WKT/EWKT geometry representation
func (p *MultiPoint) String() string {
	var s string
	if p.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", p.srid)
	}
	s += "MULTIPOINT"
	if !p.HasZ() && p.HasM() {
		s += "M"
	}

	return s + printMultiPoint(p, p.HasZ(), p.HasM())
}

// Scan implements sql.Scanner interface
func (p *MultiPoint) Scan(src interface{}) error {
	return scanGeometry(src, p)
}

// Value implements sql driver.Valuer interface
func (p *MultiPoint) Value() (driver.Value, error) {
	return p.String(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface
func (p *MultiPoint) UnmarshalBinary(data []byte) error {
	h, byteOrder, offset := readHeader(data)
	if h.Type() != MultiPointType {
		return errors.New("not expected geometry type")
	}

	p.header = h

	var err error
	p.mp, _, err = readMultiPoint(data[offset:], byteOrder, getReadPointFunc(h.wkbType))
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler interface
func (p *MultiPoint) MarshalBinary() ([]byte, error) {
	size := headerSize(p.HasSRID()) + multiPointSize(p, p.HasZ(), p.HasM())
	b := make([]byte, size)

	byteOrder := getBinaryByteOrder(p.ByteOrder())
	offset := writeHeader(p, p.Type(), byteOrder, p.HasSRID(), b)
	writeMultiPoint(p, byteOrder, p.HasZ(), p.HasM(), b[offset:])

	return b, nil
}

func printMultiPoint(p geo.MultiPoint, hasZ, hasM bool) string {
	if p.Len() == 0 {
		return " EMPTY"
	}

	var s string
	for idx := 0; idx < p.Len(); idx++ {
		s += printPoint(p.Point(idx), hasZ, hasM, false) + ","
	}

	return "(" + s[:len(s)-1] + ")"
}
