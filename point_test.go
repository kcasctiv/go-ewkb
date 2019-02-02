package ewkb

import (
	"testing"

	"github.com/kcasctiv/go-ewkb/geo"
)

func TestNewPoint(t *testing.T) {
	cases := []struct {
		name  string
		base  Base
		point geo.Point
	}{
		{"little endian", NewBase(NDR, false, false, false, 0), geo.NewPoint(6, 5)},
		{"big endian", NewBase(XDR, false, false, false, 0), geo.NewPoint(3, 1)},
		{"has Z", NewBase(NDR, true, false, false, 0), geo.NewPointZ(4, 7, 9)},
		{"has M", NewBase(NDR, false, true, false, 0), geo.NewPointM(10, 12, 17)},
		{"has Z and M", NewBase(NDR, true, true, false, 0), geo.NewPointZM(10, 12, 17, 21)},
		{"has SRID", NewBase(NDR, false, false, true, 432), geo.NewPoint(27, 43)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := NewPoint(c.base, c.point)

			if typ := p.Type(); typ != PointType {
				t.Errorf("Type: expected %v, got %v\n", PointType, typ)
			}

			if byteOrder := p.ByteOrder(); byteOrder != c.base.ByteOrder() {
				t.Errorf("ByteOrder: expected %v, got %v\n", c.base.ByteOrder(), byteOrder)
			}

			if hasZ := p.HasZ(); hasZ != c.base.HasZ() {
				t.Errorf("HasZ: expected %v, got %v\n", c.base.HasZ(), hasZ)
			}

			if hasM := p.HasM(); hasM != c.base.HasM() {
				t.Errorf("HasM: expected %v, got %v\n", c.base.HasM(), hasM)
			}

			if hasSRID := p.HasSRID(); hasSRID != c.base.HasSRID() {
				t.Errorf("HasSRID: expected %v, got %v\n", c.base.HasSRID(), hasSRID)
			}

			if srid := p.SRID(); srid != c.base.SRID() {
				t.Errorf("SRID: expected %v, got %v\n", c.base.SRID(), srid)
			}

			if x := p.X(); x != c.point.X() {
				t.Errorf("X: expected %v, got %v\n", c.point.X(), x)
			}

			if y := p.Y(); y != c.point.Y() {
				t.Errorf("Y: expected %v, got %v\n", c.point.Y(), y)
			}

			if z := p.Z(); z != c.point.Z() {
				t.Errorf("Z: expected %v, got %v\n", c.point.Z(), z)
			}

			if m := p.M(); m != c.point.M() {
				t.Errorf("M: expected %v, got %v\n", c.point.M(), m)
			}
		})
	}
}
