package ewkb

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/hex"
	"errors"
	"fmt"
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

// Available types of geometry objects
const (
	PointType uint32 = 1 + iota
	LineType
	PolygonType
	MultiPointType
	MultiLineType
	MultiPolygonType
)

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

type bbox struct {
	xmin, xmax float64
	ymin, ymax float64
	zmin, zmax float64
	mmin, mmax float64
}

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
