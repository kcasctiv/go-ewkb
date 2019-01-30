package ewkb

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

// GeometryCollection presents collection of geometry objects
type GeometryCollection struct {
	header
	geoms []Geometry
}

// NewGeometryCollection returns new GeometryCollection,
// created from geometry base and coords data
func NewGeometryCollection(b Base, geoms []Geometry) GeometryCollection {
	return GeometryCollection{
		header: header{
			byteOrder: b.ByteOrder(),
			wkbType: getFlags(
				b.HasZ(),
				b.HasM(),
				b.HasSRID(),
			) | CollectionType,
			srid: b.SRID(),
		},
		geoms: geoms,
	}
}

// Geometry returns geometry with specified index
func (c *GeometryCollection) Geometry(idx int) Geometry { return c.geoms[idx] }

// Len returns length of collection (count of geometry objects)
func (c *GeometryCollection) Len() int { return len(c.geoms) }

// String returns WKT/EWKT geometry representation
func (c *GeometryCollection) String() string {
	var s string
	if c.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", c.srid)
	}
	s += "GEOMETRYCOLLECTION"
	if !c.HasZ() && c.HasM() {
		s += "M"
	}

	if c.Len() == 0 {
		s += " EMPTY"
		return s
	}

	s += "("
	for idx := 0; idx < c.Len(); idx++ {
		gs := c.Geometry(idx).String()
		if c.Geometry(idx).HasSRID() {
			splitted := strings.Split(gs, ";")
			if len(splitted) > 1 {
				gs = splitted[1]
			}
		}

		s += gs + ","
	}

	return s[:len(s)-1] + ")"
}

// Scan implements sql.Scanner interface
func (c *GeometryCollection) Scan(src interface{}) error {
	return scanGeometry(src, c)
}

// Value implements sql driver.Valuer interface
func (c *GeometryCollection) Value() (driver.Value, error) {
	return c.String(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface
func (c *GeometryCollection) UnmarshalBinary(data []byte) error {
	h, byteOrder, offset := readHeader(data)
	if h.Type() != CollectionType {
		return errors.New("not expected geometry type")
	}

	geoms, _, err := readCollection(data[offset:], byteOrder)
	if err != nil {
		return err
	}

	c.header = h
	c.geoms = geoms

	return nil
}
