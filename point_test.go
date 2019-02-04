package ewkb

import (
	"bytes"
	"math"
	"testing"
	"time"

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

func TestPoint_UnmarshalBinary(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		valid    bool
		expected Point
	}{
		{
			"simple",
			[]byte{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 28, 64, 0, 0, 0, 0, 0, 0, 32, 64},
			true,
			NewPoint(NewBase(NDR, false, false, false, 0), geo.NewPoint(7, 8)),
		},
		{
			"with Z dimension",
			[]byte{
				1, 1, 0, 0, 128, 102, 102, 102, 102, 102, 38, 70, 192,
				205, 204, 204, 204, 204, 12, 78, 64, 0, 0, 0, 0, 0, 160, 69, 64,
			},
			true,
			NewPoint(NewBase(NDR, true, false, false, 0), geo.NewPointZ(-44.3, 60.1, 43.25)),
		},
		{
			"with M dimension",
			[]byte{
				1, 1, 0, 0, 64, 0, 0, 0, 0, 0, 0, 28, 64, 0,
				0, 0, 0, 0, 0, 32, 64, 0, 0, 0, 0, 0, 0, 34, 64,
			},
			true,
			NewPoint(NewBase(NDR, false, true, false, 0), geo.NewPointM(7, 8, 9)),
		},
		{
			"with Z and M dimension",
			[]byte{
				1, 1, 0, 0, 192, 0, 0, 0, 0, 0, 0, 240,
				63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0,
				0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 16, 64,
			},
			true,
			NewPoint(NewBase(NDR, true, true, false, 0), geo.NewPointZM(1, 2, 3, 4)),
		},
		{
			"with SRID",
			[]byte{
				1, 1, 0, 0, 32, 230, 16, 0, 0, 102, 102, 102, 102,
				102, 38, 70, 192, 205, 204, 204, 204, 204, 12, 78, 64,
			},
			true,
			NewPoint(NewBase(NDR, false, false, true, 4326), geo.NewPoint(-44.3, 60.1)),
		},
		{
			"not a point",
			[]byte{1, 2, 0, 0, 0, 0, 0, 0, 0},
			false,
			Point{},
		},
		{
			"simple corrupted",
			[]byte{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 28, 64, 0, 0, 0, 0, 0, 0, 32},
			false,
			Point{},
		},
		{
			"with Z dimension corrupted",
			[]byte{
				1, 1, 0, 0, 128, 102, 102, 102, 102, 102, 38, 70, 192,
				205, 204, 204, 204, 204, 12, 78, 64, 0, 0, 0, 0, 0, 160, 69,
			},
			false,
			Point{},
		},
		{
			"with M dimension corrupted",
			[]byte{
				1, 1, 0, 0, 64, 0, 0, 0, 0, 0, 0, 28, 64, 0,
				0, 0, 0, 0, 0, 32, 64, 0, 0, 0, 0, 0, 0, 34,
			},
			false,
			Point{},
		},
		{
			"with Z and M dimension corrupted",
			[]byte{
				1, 1, 0, 0, 192, 0, 0, 0, 0, 0, 0, 240,
				63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0,
				0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 16,
			},
			false,
			Point{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var p Point
			err := p.UnmarshalBinary(c.data)
			if err != nil && c.valid {
				t.Fatalf("Expected: no errors, got error: %v\n", err)
			}
			if err == nil && !c.valid {
				t.Fatal("Expected: error, got: no errors\n")
			}
			if !c.valid {
				return
			}

			if typ := p.Type(); typ != PointType {
				t.Errorf("Type: expected %v, got %v\n", PointType, typ)
			}

			if byteOrder := p.ByteOrder(); byteOrder != c.expected.ByteOrder() {
				t.Errorf("ByteOrder: expected %v, got %v\n", c.expected.ByteOrder(), byteOrder)
			}

			if hasZ := p.HasZ(); hasZ != c.expected.HasZ() {
				t.Errorf("HasZ: expected %v, got %v\n", c.expected.HasZ(), hasZ)
			}

			if hasM := p.HasM(); hasM != c.expected.HasM() {
				t.Errorf("HasM: expected %v, got %v\n", c.expected.HasM(), hasM)
			}

			if hasSRID := p.HasSRID(); hasSRID != c.expected.HasSRID() {
				t.Errorf("HasSRID: expected %v, got %v\n", c.expected.HasSRID(), hasSRID)
			}

			if srid := p.SRID(); srid != c.expected.SRID() {
				t.Errorf("SRID: expected %v, got %v\n", c.expected.SRID(), srid)
			}

			if x := p.X(); x != c.expected.X() {
				t.Errorf("X: expected %v, got %v\n", c.expected.X(), x)
			}

			if y := p.Y(); y != c.expected.Y() {
				t.Errorf("Y: expected %v, got %v\n", c.expected.Y(), y)
			}

			if z := p.Z(); z != c.expected.Z() {
				t.Errorf("Z: expected %v, got %v\n", c.expected.Z(), z)
			}

			if m := p.M(); m != c.expected.M() {
				t.Errorf("M: expected %v, got %v\n", c.expected.M(), m)
			}
		})
	}
}

func TestPoint_Scan(t *testing.T) {
	cases := []struct {
		name     string
		src      interface{}
		valid    bool
		expected Point
	}{
		{
			"binary",
			[]byte{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 28, 64, 0, 0, 0, 0, 0, 0, 32, 64},
			true,
			NewPoint(NewBase(NDR, false, false, false, 0), geo.NewPoint(7, 8)),
		},
		{
			"hex binary",
			[]byte{
				48, 49, 48, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
				48, 48, 48, 48, 48, 48, 48, 48, 49, 99, 52, 48, 48, 48,
				48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 50, 48, 52, 48,
			},
			true,
			NewPoint(NewBase(NDR, false, false, false, 0), geo.NewPoint(7, 8)),
		},
		{
			"hex string",
			"01010000000000000000001c400000000000002040",
			true,
			NewPoint(NewBase(NDR, false, false, false, 0), geo.NewPoint(7, 8)),
		},
		{
			"not valid hex string",
			"01010000000000000000001g400000000000002040",
			false,
			Point{},
		},
		{
			"not valid hex binary",
			[]byte{
				48, 49, 48, 49, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
				48, 48, 48, 48, 48, 48, 48, 48, 49, 99, 52, 48, 48, 48,
				48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 50, 148, 52, 48,
			},
			false,
			Point{},
		},
		{
			"not valid data type",
			time.Now(),
			false,
			Point{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var p Point
			err := p.Scan(c.src)
			if err != nil && c.valid {
				t.Fatalf("Expected: no errors, got error: %v\n", err)
			}
			if err == nil && !c.valid {
				t.Fatal("Expected: error, got: no errors\n")
			}
			if !c.valid {
				return
			}

			if typ := p.Type(); typ != PointType {
				t.Errorf("Type: expected %v, got %v\n", PointType, typ)
			}

			if byteOrder := p.ByteOrder(); byteOrder != c.expected.ByteOrder() {
				t.Errorf("ByteOrder: expected %v, got %v\n", c.expected.ByteOrder(), byteOrder)
			}

			if hasZ := p.HasZ(); hasZ != c.expected.HasZ() {
				t.Errorf("HasZ: expected %v, got %v\n", c.expected.HasZ(), hasZ)
			}

			if hasM := p.HasM(); hasM != c.expected.HasM() {
				t.Errorf("HasM: expected %v, got %v\n", c.expected.HasM(), hasM)
			}

			if hasSRID := p.HasSRID(); hasSRID != c.expected.HasSRID() {
				t.Errorf("HasSRID: expected %v, got %v\n", c.expected.HasSRID(), hasSRID)
			}

			if srid := p.SRID(); srid != c.expected.SRID() {
				t.Errorf("SRID: expected %v, got %v\n", c.expected.SRID(), srid)
			}

			if x := p.X(); x != c.expected.X() {
				t.Errorf("X: expected %v, got %v\n", c.expected.X(), x)
			}

			if y := p.Y(); y != c.expected.Y() {
				t.Errorf("Y: expected %v, got %v\n", c.expected.Y(), y)
			}

			if z := p.Z(); z != c.expected.Z() {
				t.Errorf("Z: expected %v, got %v\n", c.expected.Z(), z)
			}

			if m := p.M(); m != c.expected.M() {
				t.Errorf("M: expected %v, got %v\n", c.expected.M(), m)
			}
		})
	}
}

func TestPoint_String(t *testing.T) {
	cases := []struct {
		name     string
		point    Point
		expected string
	}{
		{
			"simple",
			NewPoint(NewBase(NDR, false, false, false, 0), geo.NewPoint(6, 5)),
			"POINT(6 5)",
		},
		{
			"with Z dimension",
			NewPoint(NewBase(NDR, true, false, false, 0), geo.NewPointZ(6, 5, 4)),
			"POINT(6 5 4)",
		},
		{
			"with M dimension",
			NewPoint(NewBase(NDR, false, true, false, 0), geo.NewPointM(6, 5, 4)),
			"POINTM(6 5 4)",
		},
		{
			"with Z and M dimensions",
			NewPoint(NewBase(NDR, true, true, false, 0), geo.NewPointZM(6, 5, 4, 3)),
			"POINT(6 5 4 3)",
		},
		{
			"with SRID",
			NewPoint(NewBase(NDR, false, false, true, 4321), geo.NewPoint(6, 5)),
			"SRID=4321;POINT(6 5)",
		},
		{
			"empty",
			NewPoint(NewBase(NDR, false, false, false, 0), geo.NewPoint(math.NaN(), math.NaN())),
			"POINT EMPTY",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := c.point.String()

			if s != c.expected {
				t.Errorf("Expected %q, got %q\n", c.expected, s)
			}
		})
	}
}

func TestPoint_Value(t *testing.T) {
	cases := []struct {
		name     string
		point    Point
		expected string
	}{
		{
			"simple",
			NewPoint(NewBase(NDR, false, false, false, 0), geo.NewPoint(6, 5)),
			"POINT(6 5)",
		},
		{
			"with Z dimension",
			NewPoint(NewBase(NDR, true, false, false, 0), geo.NewPointZ(6, 5, 4)),
			"POINT(6 5 4)",
		},
		{
			"with M dimension",
			NewPoint(NewBase(NDR, false, true, false, 0), geo.NewPointM(6, 5, 4)),
			"POINTM(6 5 4)",
		},
		{
			"with Z and M dimensions",
			NewPoint(NewBase(NDR, true, true, false, 0), geo.NewPointZM(6, 5, 4, 3)),
			"POINT(6 5 4 3)",
		},
		{
			"with SRID",
			NewPoint(NewBase(NDR, false, false, true, 4321), geo.NewPoint(6, 5)),
			"SRID=4321;POINT(6 5)",
		},
		{
			"empty",
			NewPoint(NewBase(NDR, false, false, false, 0), geo.NewPoint(math.NaN(), math.NaN())),
			"POINT EMPTY",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, err := c.point.Value()
			if err != nil {
				t.Fatalf("Expected not errors, got %v\n", err)
			}

			if s != c.expected {
				t.Errorf("Expected %q, got %q\n", c.expected, s)
			}
		})
	}
}

func TestPoint_MarshalBinary(t *testing.T) {
	cases := []struct {
		name     string
		data     Point
		expected []byte
	}{
		{
			"simple",
			NewPoint(NewBase(NDR, false, false, false, 0), geo.NewPoint(7, 8)),
			[]byte{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 28, 64, 0, 0, 0, 0, 0, 0, 32, 64},
		},
		{
			"with Z dimension",
			NewPoint(NewBase(NDR, true, false, false, 0), geo.NewPointZ(-44.3, 60.1, 43.25)),
			[]byte{
				1, 1, 0, 0, 128, 102, 102, 102, 102, 102, 38, 70, 192,
				205, 204, 204, 204, 204, 12, 78, 64, 0, 0, 0, 0, 0, 160, 69, 64,
			},
		},
		{
			"with M dimension",
			NewPoint(NewBase(NDR, false, true, false, 0), geo.NewPointM(7, 8, 9)),
			[]byte{
				1, 1, 0, 0, 64, 0, 0, 0, 0, 0, 0, 28, 64, 0,
				0, 0, 0, 0, 0, 32, 64, 0, 0, 0, 0, 0, 0, 34, 64,
			},
		},
		{
			"with Z and M dimension",
			NewPoint(NewBase(NDR, true, true, false, 0), geo.NewPointZM(1, 2, 3, 4)),
			[]byte{
				1, 1, 0, 0, 192, 0, 0, 0, 0, 0, 0, 240,
				63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0,
				0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 16, 64,
			},
		},
		{
			"with SRID",
			NewPoint(NewBase(NDR, false, false, true, 4326), geo.NewPoint(-44.3, 60.1)),
			[]byte{
				1, 1, 0, 0, 32, 230, 16, 0, 0, 102, 102, 102, 102,
				102, 38, 70, 192, 205, 204, 204, 204, 204, 12, 78, 64,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, err := c.data.MarshalBinary()
			if err != nil {
				t.Fatalf("Expected not errors, got %v\n", err)
			}

			if !bytes.Equal(b, c.expected) {
				t.Errorf("Expected %v, got %v\n", c.expected, b)
			}
		})
	}
}
