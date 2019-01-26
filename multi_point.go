package ewkb

import (
	"database/sql/driver"
	"fmt"

	"github.com/kcasctiv/go-ewkb/geo"
)

type MultiPoint struct {
	header
	mp geo.MultiPoint
}

func (p *MultiPoint) Point(idx int) geo.Point { return p.mp.Point(idx) }
func (p *MultiPoint) Len() int                { return p.mp.Len() }

func (p *MultiPoint) String() string {
	var s string
	if p.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", p.srid)
	}
	s += "MULTIPOINT "
	if p.HasZ() {
		s += "Z"
	}
	if p.HasM() {
		s += "M"
	}

	return s + " " + printMultiPoint(p, p.HasZ(), p.HasM())
}

func (p *MultiPoint) Scan(src interface{}) error {
	// TODO:
	return nil
}

func (p *MultiPoint) Value() (driver.Value, error) {
	return p.String(), nil
}

func printMultiPoint(p geo.MultiPoint, hasZ, hasM bool) string {
	if p.Len() == 0 {
		return "()"
	}

	var s string
	for idx := 0; idx < p.Len(); idx++ {
		s += printPoint(p.Point(idx), hasZ, hasM) + ", "
	}

	return "(" + s[:len(s)-2] + ")"
}
