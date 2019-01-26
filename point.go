package ewkb

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/kcasctiv/go-ewkb/geo"
)

// Point presents 2, 3 or 4 dimensions point
type Point struct {
	header
	point geo.Point
}

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
	h, byteOrder, offset := readHeader(data)
	if h.Type() != PointType {
		return errors.New("not expected geometry type")
	}

	p.header = h

	var err error
	p.point, _, err = getReadPointFunc(h.wkbType)(data[offset:], byteOrder)
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
