package ewkb

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/kcasctiv/go-ewkb/geo"
)

// LineString presents LineString geometry object
type LineString struct {
	header
	mp geo.MultiPoint
}

// Point returns point of LineString with specified index
func (l *LineString) Point(idx int) geo.Point { return l.mp.Point(idx) }

// Len returns length of LineString (count of points)
func (l *LineString) Len() int { return l.mp.Len() }

// String returns WKT/EWKT geometry representation
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

// Scan implements sql.Scanner interface
func (l *LineString) Scan(src interface{}) error {
	return scanGeometry(src, l)
}

// Value implements sql driver.Valuer interface
func (l *LineString) Value() (driver.Value, error) {
	return l.String(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface
func (l *LineString) UnmarshalBinary(data []byte) error {
	h, byteOrder, offset := readHeader(data)
	if h.Type() != LineType {
		return errors.New("not expected geometry type")
	}

	l.header = h

	var err error
	l.mp, _, err = readMultiPoint(data[offset:], byteOrder, getReadPointFunc(h.wkbType))
	return err
}