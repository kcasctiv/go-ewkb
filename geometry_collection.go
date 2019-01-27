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
	s += "GEOMETRYCOLLECTION "
	if c.HasZ() {
		s += "Z"
	}
	if c.HasM() {
		s += "M"
	}

	if c.Len() == 0 {
		s += "EMPTY"
		return s
	}

	s += " ("
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

	c.header = h

	length := byteOrder.Uint32(data[offset:])
	offset += 4

	geoms := make([]Geometry, length)
	var n int
	var err error
	for idx := 0; idx < len(geoms); idx++ {
		h1, byteOrder1, offset1 := readHeader(data[offset:])
		offset += offset1
		switch h1.Type() {
		case PointType:
			point := Point{header: h1}
			point.point, n, err = getReadPointFunc(h1.wkbType)(data[offset:], byteOrder1)
			geoms[idx] = &point
		case LineType:
			line := LineString{header: h1}
			line.mp, n, err = readMultiPoint(data[offset:], byteOrder1, getReadPointFunc(h1.wkbType))
			geoms[idx] = &line
		case PolygonType:
			poly := Polygon{header: h1}
			poly.poly, n, err = readPolygon(data[offset:], byteOrder1, getReadPointFunc(h1.wkbType))
			geoms[idx] = &poly
		case MultiPointType:
			mpoint := MultiPoint{header: h1}
			mpoint.mp, n, err = readMultiPoint(data[offset:], byteOrder1, getReadPointFunc(h1.wkbType))
			geoms[idx] = &mpoint
		case MultiLineType:
			mline := MultiLineString{header: h1}
			mline.ml, n, err = readMultiLine(data[offset:], byteOrder1, getReadPointFunc(h1.wkbType))
			geoms[idx] = &mline
		case MultiPolygonType:
			mpoly := MultiPolygon{header: h1}
			mpoly.mp, n, err = readMultiPolygon(data[offset:], byteOrder1, getReadPointFunc(h1.wkbType))
			geoms[idx] = &mpoly
		default:
			return errors.New("not expected geometry type")
		}

		if err != nil {
			return err
		}
		offset += n
	}

	c.geoms = geoms

	return nil
}
