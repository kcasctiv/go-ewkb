package ewkb

import (
	"database/sql/driver"
	"fmt"
	"math"

	"github.com/kcasctiv/go-ewkb/geo"
)

type MultiPolygon struct {
	byteOrder byte
	wkbType   uint32
	srid      int32
	bbox      *bbox
	mp        geo.MultiPolygon
}

func (p *MultiPolygon) ByteOrder() byte             { return p.byteOrder }
func (p *MultiPolygon) Type() uint32                { return p.wkbType & uint32(math.MaxUint16) }
func (p *MultiPolygon) HasZ() bool                  { return (p.wkbType & zFlag) == zFlag }
func (p *MultiPolygon) HasM() bool                  { return (p.wkbType & mFlag) == mFlag }
func (p *MultiPolygon) HasSRID() bool               { return (p.wkbType & sridFlag) == sridFlag }
func (p *MultiPolygon) HasBBOX() bool               { return (p.wkbType & bboxFlag) == bboxFlag }
func (p *MultiPolygon) Polygon(idx int) geo.Polygon { return p.mp.Polygon(idx) }
func (p *MultiPolygon) Len() int                    { return p.mp.Len() }

func (p *MultiPolygon) String() string {
	var s string
	if p.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", p.srid)
	}
	s += "MULTIPOLYGON "
	if p.HasZ() {
		s += "Z"
	}
	if p.HasM() {
		s += "M"
	}
	s += " ("
	if p.Len() > 0 {
		for idx := 0; idx < p.Len(); idx++ {
			s += printPolygon(p.Polygon(idx), p.HasZ(), p.HasM()) + ", "
		}

		s = s[:len(s)-2]
	}

	return s + ")"
}

func (p *MultiPolygon) Scan(src interface{}) error {
	// TODO:
	return nil
}

func (p *MultiPolygon) Value() (driver.Value, error) {
	return p.String(), nil
}
