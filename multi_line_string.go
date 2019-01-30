package ewkb

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/kcasctiv/go-ewkb/geo"
)

// MultiLineString presents MultiLineString geometry object
type MultiLineString struct {
	header
	ml geo.MultiLine
}

// NewMultiLineString returns new MultiLineString,
// created from geometry base and coords data
func NewMultiLineString(b Base, ml geo.MultiLine) MultiLineString {
	return MultiLineString{
		header: header{
			byteOrder: b.ByteOrder(),
			wkbType: getFlags(
				b.HasZ(),
				b.HasM(),
				b.HasSRID(),
			) | MultiLineType,
			srid: b.SRID(),
		},
		ml: ml,
	}
}

// Line returns line with specified index
func (l *MultiLineString) Line(idx int) geo.MultiPoint { return l.ml.Line(idx) }

// Len returns count of lines
func (l *MultiLineString) Len() int { return l.ml.Len() }

// String returns WKT/EWKT geometry representation
func (l *MultiLineString) String() string {
	var s string
	if l.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", l.srid)
	}
	s += "MULTILINESTRING"
	if !l.HasZ() && l.HasM() {
		s += "M"
	}

	if l.Len() == 0 {
		s += " EMPTY"
		return s
	}

	s += "("
	if l.Len() > 0 {
		for idx := 0; idx < l.Len(); idx++ {
			s += printMultiPoint(l.Line(idx), l.HasZ(), l.HasM()) + ","
		}

		s = s[:len(s)-1]
	}

	return s + ")"
}

// Scan implements sql.Scanner interface
func (l *MultiLineString) Scan(src interface{}) error {
	return scanGeometry(src, l)
}

// Value implements sql driver.Valuer interface
func (l *MultiLineString) Value() (driver.Value, error) {
	return l.String(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface
func (l *MultiLineString) UnmarshalBinary(data []byte) error {
	h, byteOrder, offset := readHeader(data)
	if h.Type() != MultiLineType {
		return errors.New("not expected geometry type")
	}

	l.header = h

	var err error
	l.ml, _, err = readMultiLine(
		data[offset:],
		byteOrder,
		getReadPointFunc(h.wkbType),
	)
	return err
}
