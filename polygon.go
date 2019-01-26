package ewkb

import (
	"database/sql/driver"
	"fmt"

	"github.com/kcasctiv/go-ewkb/geo"
)

type Polygon struct {
	header
	poly geo.Polygon
}

func (p *Polygon) Ring(idx int) geo.MultiPoint { return p.poly.Ring(idx) }
func (p *Polygon) Len() int                    { return p.poly.Len() }

func (p *Polygon) String() string {
	var s string
	if p.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", p.srid)
	}
	s += "POLYGON "
	if p.HasZ() {
		s += "Z"
	}
	if p.HasM() {
		s += "M"
	}

	return s + " " + printPolygon(p, p.HasZ(), p.HasM())
}

func (p *Polygon) Scan(src interface{}) error {
	// TODO:
	return nil
}

func (p *Polygon) Value() (driver.Value, error) {
	return p.String(), nil
}

func printPolygon(p geo.Polygon, hasZ, hasM bool) string {
	if p.Len() == 0 {
		return "()"
	}

	var s string
	for idx := 0; idx < p.Len(); idx++ {
		s += printMultiPoint(p.Ring(idx), hasZ, hasM) + ", "
	}

	return "(" + s[:len(s)-2] + ")"
}
