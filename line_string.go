package ewkb

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/kcasctiv/go-ewkb/geo"
)

// LineString presents LineString geometry object
type LineString struct {
	byteOrder byte
	wkbType   uint32
	srid      int32
	bbox      *bbox
	mp        geo.MultiPoint
}

// ByteOrder returns byte order of geometry
func (l *LineString) ByteOrder() byte { return l.byteOrder }

// Type returns type of geometry
func (l *LineString) Type() uint32 { return l.wkbType & uint32(math.MaxUint16) }

// HasZ checks if geometry has Z dimension
func (l *LineString) HasZ() bool { return (l.wkbType & zFlag) == zFlag }

// HasM checks if geometry has M dimension
func (l *LineString) HasM() bool { return (l.wkbType & mFlag) == mFlag }

// HasSRID checks if geometry contains SRID
func (l *LineString) HasSRID() bool { return (l.wkbType & sridFlag) == sridFlag }

// HasBBOX checks if geometry contains BBOX
func (l *LineString) HasBBOX() bool { return (l.wkbType & bboxFlag) == bboxFlag }

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
	var byteOrder binary.ByteOrder
	if data[0] == XDR {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}

	offset := 1
	wkbType := byteOrder.Uint32(data[offset:])
	if (wkbType & uint32(math.MaxUint16)) != LineType {
		return errors.New("not expected geometry type")
	}
	l.byteOrder = data[0]
	l.wkbType = wkbType
	offset += 4

	if (wkbType & sridFlag) == sridFlag {
		l.srid = int32(byteOrder.Uint32(data[offset:]))
		offset += 4
	}

	if (wkbType & bboxFlag) == bboxFlag {
		// TODO:
	}

	var err error
	l.mp, _, err = readMultiPoint(data[offset:], byteOrder, getReadPointFunc(wkbType))
	return err
}
