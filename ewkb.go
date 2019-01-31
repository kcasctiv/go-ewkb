package ewkb

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
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
)

// Available types of geometry objects
const (
	PointType uint32 = 1 + iota
	LineType
	PolygonType
	MultiPointType
	MultiLineType
	MultiPolygonType
	CollectionType
)

// Base presents interface of base of geometry
type Base interface {
	// ByteOrder returns byte order of geometry
	ByteOrder() byte
	// HasZ checks if geometry has Z dimension
	HasZ() bool
	// HasM checks if geometry has M dimension
	HasM() bool
	// HasSRID checks if geometry contains SRID
	HasSRID() bool
	// SRID returns SRID, or zero, if there is no SRID
	SRID() int32
}

// Geometry presents interface of geometry object
type Geometry interface {
	// Type returns type of geometry
	Type() uint32
	Base
	fmt.Stringer
	sql.Scanner
	driver.Valuer
	encoding.BinaryUnmarshaler
	encoding.BinaryMarshaler
}

// NewBase returns new base of geometry
func NewBase(byteOrder byte, hasZ, hasM, hasSRID bool, srid int32) Base {
	return &header{
		byteOrder: byteOrder,
		wkbType:   getFlags(hasZ, hasM, hasSRID),
		srid:      srid,
	}
}

// Wrapper prensents wrapper for geometry objects.
// Can be used for reading from and writing to DB
// all types of geometry, supported by package.
// Also support null values and may be useful
// for nullable columns
type Wrapper struct {
	Geometry Geometry
}

// Scan implements sql.Scanner interface
func (w *Wrapper) Scan(src interface{}) error {
	if src == nil {
		w.Geometry = nil
		return nil
	}

	return scanGeometry(src, w)
}

// Value implements sql driver.Valuer interface
func (w *Wrapper) Value() (driver.Value, error) {
	if w.Geometry == nil {
		return nil, nil
	}

	return w.Geometry.Value()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface
func (w *Wrapper) UnmarshalBinary(data []byte) error {
	if data == nil {
		w.Geometry = nil
		return nil
	}

	var err error
	h, byteOrder, offset := readHeader(data)
	switch h.Type() {
	case PointType:
		point := Point{header: h}
		point.point, _, err = getReadPointFunc(h.wkbType)(data[offset:], byteOrder)
		w.Geometry = &point
	case LineType:
		line := LineString{header: h}
		line.mp, _, err = readMultiPoint(data[offset:], byteOrder, getReadPointFunc(h.wkbType))
		w.Geometry = &line
	case PolygonType:
		poly := Polygon{header: h}
		poly.poly, _, err = readPolygon(data[offset:], byteOrder, getReadPointFunc(h.wkbType))
		w.Geometry = &poly
	case MultiPointType:
		mpoint := MultiPoint{header: h}
		mpoint.mp, _, err = readMultiPoint(data[offset:], byteOrder, getReadPointFunc(h.wkbType))
		w.Geometry = &mpoint
	case MultiLineType:
		mline := MultiLineString{header: h}
		mline.ml, _, err = readMultiLine(data[offset:], byteOrder, getReadPointFunc(h.wkbType))
		w.Geometry = &mline
	case MultiPolygonType:
		mpoly := MultiPolygon{header: h}
		mpoly.mp, _, err = readMultiPolygon(data[offset:], byteOrder, getReadPointFunc(h.wkbType))
		w.Geometry = &mpoly
	case CollectionType:
		gc := GeometryCollection{header: h}
		gc.geoms, _, err = readCollection(data[offset:], byteOrder)
		w.Geometry = &gc
	default:
		err = errors.New("not expected geometry type")
	}
	if err != nil {
		w.Geometry = nil
	}

	return err
}

// MarshalBinary implements encoding.BinaryMarshaler interface
func (w *Wrapper) MarshalBinary() ([]byte, error) {
	if w.Geometry == nil {
		return nil, nil
	}

	return w.Geometry.MarshalBinary()
}

type header struct {
	byteOrder byte
	wkbType   uint32
	srid      int32
}

// ByteOrder returns byte order of geometry
func (h *header) ByteOrder() byte { return h.byteOrder }

// Type returns type of geometry
func (h *header) Type() uint32 { return h.wkbType & uint32(math.MaxUint16) }

// HasZ checks if geometry has Z dimension
func (h *header) HasZ() bool { return (h.wkbType & zFlag) == zFlag }

// HasM checks if geometry has M dimension
func (h *header) HasM() bool { return (h.wkbType & mFlag) == mFlag }

// HasSRID checks if geometry contains SRID
func (h *header) HasSRID() bool { return (h.wkbType & sridFlag) == sridFlag }

// SRID returns SRID, or zero, if there is no SRID
func (h *header) SRID() int32 { return h.srid }

func scanGeometry(src interface{}, unmarshaler encoding.BinaryUnmarshaler) error {
	var data []byte
	var err error
	switch d := src.(type) {
	case []byte:
		if d[0] == 48 {
			data, err = hex.DecodeString(string(d))
		} else {
			data = d
		}
	case string:
		data, err = hex.DecodeString(d)
	default:
		return errors.New("could not scan geometry")
	}
	if err != nil {
		return fmt.Errorf("could not scan geometry: %v", err)
	}

	return unmarshaler.UnmarshalBinary(data)
}

func getFlags(z, m, srid bool) uint32 {
	var flags uint32
	if z {
		flags = flags | zFlag
	}

	if m {
		flags = flags | mFlag
	}

	if srid {
		flags = flags | sridFlag
	}

	return flags
}

func getBinaryByteOrder(b byte) binary.ByteOrder {
	if b == XDR {
		return binary.BigEndian
	}

	return binary.LittleEndian
}
