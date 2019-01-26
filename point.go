package ewkb

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/kcasctiv/go-ewkb/geo"
)

// Point presents 2, 3 or 4 dimensions point
type Point struct {
	byteOrder byte
	wkbType   uint32
	srid      int32
	bbox      *bbox
	point     geo.Point
}

// ByteOrder returns byte order of geometry
func (p *Point) ByteOrder() byte { return p.byteOrder }

// Type returns type of geometry
func (p *Point) Type() uint32 { return p.wkbType & uint32(math.MaxUint16) }

// HasZ checks if geometry has Z dimension
func (p *Point) HasZ() bool { return (p.wkbType & zFlag) == zFlag }

// HasM checks if geometry has M dimension
func (p *Point) HasM() bool { return (p.wkbType & mFlag) == mFlag }

// HasSRID checks if geometry contains SRID
func (p *Point) HasSRID() bool { return (p.wkbType & sridFlag) == sridFlag }

// HasBBOX checks if geometry contains BBOX
func (p *Point) HasBBOX() bool { return (p.wkbType & bboxFlag) == bboxFlag }

// X returns value of X dimension
func (p *Point) X() float64 { return p.point.X() }

// Y returns value of Y dimension
func (p *Point) Y() float64 { return p.point.Y() }

// Z returns value of Z dimension
func (p *Point) Z() float64 { return p.point.Z() }

// M returns value of M dimension
func (p *Point) M() float64 { return p.point.M() }

// String returns WKT/EWKT geometry representation
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

// Scan implements sql.Scanner interface
func (p *Point) Scan(src interface{}) error {
	return scanGeometry(src, p)
}

// Value implements sql driver.Valuer interface
func (p *Point) Value() (driver.Value, error) {
	return p.String(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface
func (p *Point) UnmarshalBinary(data []byte) error {
	var byteOrder binary.ByteOrder
	if data[0] == XDR {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}

	offset := 1
	wkbType := byteOrder.Uint32(data[offset:])
	if (wkbType & uint32(math.MaxUint16)) != PointType {
		return errors.New("not expected geometry type")
	}
	p.byteOrder = data[0]
	p.wkbType = wkbType
	offset += 4

	if (wkbType & sridFlag) == sridFlag {
		p.srid = int32(byteOrder.Uint32(data[offset:]))
		offset += 4
	}

	if (wkbType & bboxFlag) == bboxFlag {
		// TODO:
	}

	var err error
	p.point, _, err = getReadPointFunc(wkbType)(data[offset:], byteOrder)
	return err
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
