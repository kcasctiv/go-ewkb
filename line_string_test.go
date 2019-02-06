package ewkb

import (
	"bytes"
	"testing"
	"time"

	"github.com/kcasctiv/go-ewkb/geo"
)

func TestNewLineString(t *testing.T) {
	cases := []struct {
		name string
		base Base
		line geo.MultiPoint
	}{
		{
			"little endian",
			NewBase(NDR, false, false, false, 0),
			geo.NewMultiPoint([]geo.Point{
				geo.NewPoint(6, 5),
				geo.NewPoint(7, 8),
			}),
		},
		{
			"big endian",
			NewBase(XDR, false, false, false, 0),
			geo.NewMultiPoint([]geo.Point{
				geo.NewPoint(6, 5),
				geo.NewPoint(7, 8),
			}),
		},
		{
			"has Z",
			NewBase(NDR, true, false, false, 0),
			geo.NewMultiPoint([]geo.Point{
				geo.NewPointZ(4, 7, 9),
				geo.NewPointZ(6, 8, 12),
			}),
		},
		{
			"has M",
			NewBase(NDR, false, true, false, 0),
			geo.NewMultiPoint([]geo.Point{
				geo.NewPointM(4, 7, 9),
				geo.NewPointM(6, 8, 12),
			}),
		},
		{
			"has Z and M",
			NewBase(NDR, true, true, false, 0),
			geo.NewMultiPoint([]geo.Point{
				geo.NewPointZM(4, 7, 9, 14),
				geo.NewPointZM(6, 8, 12, 7),
			}),
		},
		{
			"has SRID",
			NewBase(NDR, false, false, true, 432),
			geo.NewMultiPoint([]geo.Point{
				geo.NewPoint(6, 5),
				geo.NewPoint(7, 8),
			}),
		},
		{
			"has 3 points",
			NewBase(NDR, false, false, false, 0),
			geo.NewMultiPoint([]geo.Point{
				geo.NewPoint(6, 5),
				geo.NewPoint(7, 8),
				geo.NewPoint(3, 4),
			}),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			line := NewLineString(c.base, c.line)

			if typ := line.Type(); typ != LineType {
				t.Errorf("Type: expected %v, got %v\n", LineType, typ)
			}

			if byteOrder := line.ByteOrder(); byteOrder != c.base.ByteOrder() {
				t.Errorf("ByteOrder: expected %v, got %v\n", c.base.ByteOrder(), byteOrder)
			}

			if hasZ := line.HasZ(); hasZ != c.base.HasZ() {
				t.Errorf("HasZ: expected %v, got %v\n", c.base.HasZ(), hasZ)
			}

			if hasM := line.HasM(); hasM != c.base.HasM() {
				t.Errorf("HasM: expected %v, got %v\n", c.base.HasM(), hasM)
			}

			if hasSRID := line.HasSRID(); hasSRID != c.base.HasSRID() {
				t.Errorf("HasSRID: expected %v, got %v\n", c.base.HasSRID(), hasSRID)
			}

			if srid := line.SRID(); srid != c.base.SRID() {
				t.Errorf("SRID: expected %v, got %v\n", c.base.SRID(), srid)
			}

			if len := line.Len(); len != c.line.Len() {
				t.Errorf("Len: expected %v, got %v\n", c.line.Len(), len)
			}

			for idx := 0; idx < line.Len(); idx++ {
				point := line.Point(idx)
				epoint := c.line.Point(idx)
				if x := point.X(); x != epoint.X() {
					t.Errorf("Point: %d: X: expected %v, got %v\n", idx, epoint.X(), x)
				}

				if y := point.Y(); y != epoint.Y() {
					t.Errorf("Point: %d: Y: expected %v, got %v\n", idx, epoint.Y(), y)
				}

				if z := point.Z(); z != epoint.Z() {
					t.Errorf("Point: %d: Z: expected %v, got %v\n", idx, epoint.Z(), z)
				}

				if m := point.M(); m != epoint.M() {
					t.Errorf("Point: %d: M: expected %v, got %v\n", idx, epoint.M(), m)
				}
			}
		})
	}
}

func TestLineString_UnmarshalBinary(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		valid    bool
		expected LineString
	}{
		{
			"simple",
			[]byte{
				1, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0,
				0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64,
			},
			true,
			NewLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(1, 3),
					geo.NewPoint(2, 4),
				}),
			),
		},
		{
			"with Z dimension",
			[]byte{
				1, 2, 0, 0, 128, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0,
				0, 28, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0,
				0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 20, 64,
			},
			true,
			NewLineString(
				NewBase(NDR, true, false, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointZ(1, 3, 7),
					geo.NewPointZ(2, 4, 5),
				}),
			),
		},
		{
			"with M dimension",
			[]byte{
				1, 2, 0, 0, 64, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0,
				0, 28, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0,
				0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 20, 64,
			},
			true,
			NewLineString(
				NewBase(NDR, false, true, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointM(1, 3, 7),
					geo.NewPointM(2, 4, 5),
				}),
			),
		},
		{
			"with Z and M dimension",
			[]byte{
				1, 2, 0, 0, 192, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0,
				0, 0, 28, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0,
				0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0,
				0, 0, 0, 0, 20, 64, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			true,
			NewLineString(
				NewBase(NDR, true, true, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointZM(1, 3, 7, 2),
					geo.NewPointZM(2, 4, 5, 0),
				}),
			),
		},
		{
			"with SRID",
			[]byte{
				1, 2, 0, 0, 32, 230, 16, 0, 0, 2,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 240, 63, 0,
				0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0,
				0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64,
			},
			true,
			NewLineString(
				NewBase(NDR, false, false, true, 4326),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(1, 3),
					geo.NewPoint(2, 4),
				}),
			),
		},
		{
			"not a line",
			[]byte{1, 3, 0, 0, 0, 0, 0, 0, 0},
			false,
			LineString{},
		},
		{
			"simple corrupted",
			[]byte{
				1, 2, 0, 0, 0, 2, 0, 0,
			},
			false,
			LineString{},
		},
		{
			"with Z dimension corrupted",
			[]byte{
				1, 2, 0, 0, 128, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0,
				0, 28, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0,
				0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 20, 64,
			},
			false,
			LineString{},
		},
		{
			"with M dimension corrupted",
			[]byte{
				1, 2, 0, 0, 64, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0,
				0, 28, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0,
				0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 20, 64,
			},
			false,
			LineString{},
		},
		{
			"with Z and M dimension corrupted",
			[]byte{
				1, 2, 0, 0, 192, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0,
				0, 0, 28, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0,
				0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 0, 0,
				0, 0, 0, 0, 20, 64, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			false,
			LineString{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var line LineString
			err := line.UnmarshalBinary(c.data)
			if err != nil && c.valid {
				t.Fatalf("Expected: no errors, got error: %v\n", err)
			}
			if err == nil && !c.valid {
				t.Fatal("Expected: error, got: no errors\n")
			}
			if !c.valid {
				return
			}

			if typ := line.Type(); typ != LineType {
				t.Errorf("Type: expected %v, got %v\n", LineType, typ)
			}

			if byteOrder := line.ByteOrder(); byteOrder != c.expected.ByteOrder() {
				t.Errorf("ByteOrder: expected %v, got %v\n", c.expected.ByteOrder(), byteOrder)
			}

			if hasZ := line.HasZ(); hasZ != c.expected.HasZ() {
				t.Errorf("HasZ: expected %v, got %v\n", c.expected.HasZ(), hasZ)
			}

			if hasM := line.HasM(); hasM != c.expected.HasM() {
				t.Errorf("HasM: expected %v, got %v\n", c.expected.HasM(), hasM)
			}

			if hasSRID := line.HasSRID(); hasSRID != c.expected.HasSRID() {
				t.Errorf("HasSRID: expected %v, got %v\n", c.expected.HasSRID(), hasSRID)
			}

			if srid := line.SRID(); srid != c.expected.SRID() {
				t.Errorf("SRID: expected %v, got %v\n", c.expected.SRID(), srid)
			}

			if len := line.Len(); len != c.expected.Len() {
				t.Errorf("Len: expected %v, got %v\n", c.expected.Len(), len)
			}

			for idx := 0; idx < line.Len(); idx++ {
				point := line.Point(idx)
				epoint := c.expected.Point(idx)
				if x := point.X(); x != epoint.X() {
					t.Errorf("Point: %d: X: expected %v, got %v\n", idx, epoint.X(), x)
				}

				if y := point.Y(); y != epoint.Y() {
					t.Errorf("Point: %d: Y: expected %v, got %v\n", idx, epoint.Y(), y)
				}

				if z := point.Z(); z != epoint.Z() {
					t.Errorf("Point: %d: Z: expected %v, got %v\n", idx, epoint.Z(), z)
				}

				if m := point.M(); m != epoint.M() {
					t.Errorf("Point: %d: M: expected %v, got %v\n", idx, epoint.M(), m)
				}
			}
		})
	}
}

func TestLineString_Scan(t *testing.T) {
	cases := []struct {
		name     string
		src      interface{}
		valid    bool
		expected LineString
	}{
		{
			"binary",
			[]byte{
				1, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0,
				0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64,
			},
			true,
			NewLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(1, 3),
					geo.NewPoint(2, 4),
				}),
			),
		},
		{
			"hex binary",
			[]byte{
				48, 49, 48, 50, 48, 48, 48, 48, 48, 48, 48, 50, 48,
				48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
				48, 48, 48, 48, 70, 48, 51, 70, 48, 48, 48, 48, 48,
				48, 48, 48, 48, 48, 48, 48, 48, 56, 52, 48, 48, 48,
				48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 52,
				48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
				49, 48, 52, 48,
			},
			true,
			NewLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(1, 3),
					geo.NewPoint(2, 4),
				}),
			),
		},
		{
			"hex string",
			"010200000002000000000000000000F03F000000000000084000000000000000400000000000001040",
			true,
			NewLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(1, 3),
					geo.NewPoint(2, 4),
				}),
			),
		},
		{
			"not valid hex string",
			"01020000000000000000001g400000000000002040",
			false,
			LineString{},
		},
		{
			"not valid hex binary",
			[]byte{
				48, 49, 48, 50, 48, 48, 48, 48, 48, 48, 48, 50, 48,
				48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
				48, 48, 48, 48, 70, 48, 51, 70, 48, 48, 48, 48, 48,
				48, 48, 48, 48, 154, 48, 48, 48, 56, 52, 48, 48, 48,
				48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 52,
				48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
				49, 48, 52, 48,
			},
			false,
			LineString{},
		},
		{
			"not valid data type",
			time.Now(),
			false,
			LineString{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var line LineString
			err := line.Scan(c.src)
			if err != nil && c.valid {
				t.Fatalf("Expected: no errors, got error: %v\n", err)
			}
			if err == nil && !c.valid {
				t.Fatal("Expected: error, got: no errors\n")
			}
			if !c.valid {
				return
			}

			if typ := line.Type(); typ != LineType {
				t.Errorf("Type: expected %v, got %v\n", LineType, typ)
			}

			if byteOrder := line.ByteOrder(); byteOrder != c.expected.ByteOrder() {
				t.Errorf("ByteOrder: expected %v, got %v\n", c.expected.ByteOrder(), byteOrder)
			}

			if hasZ := line.HasZ(); hasZ != c.expected.HasZ() {
				t.Errorf("HasZ: expected %v, got %v\n", c.expected.HasZ(), hasZ)
			}

			if hasM := line.HasM(); hasM != c.expected.HasM() {
				t.Errorf("HasM: expected %v, got %v\n", c.expected.HasM(), hasM)
			}

			if hasSRID := line.HasSRID(); hasSRID != c.expected.HasSRID() {
				t.Errorf("HasSRID: expected %v, got %v\n", c.expected.HasSRID(), hasSRID)
			}

			if srid := line.SRID(); srid != c.expected.SRID() {
				t.Errorf("SRID: expected %v, got %v\n", c.expected.SRID(), srid)
			}

			if len := line.Len(); len != c.expected.Len() {
				t.Errorf("Len: expected %v, got %v\n", c.expected.Len(), len)
			}

			for idx := 0; idx < line.Len(); idx++ {
				point := line.Point(idx)
				epoint := c.expected.Point(idx)
				if x := point.X(); x != epoint.X() {
					t.Errorf("Point: %d: X: expected %v, got %v\n", idx, epoint.X(), x)
				}

				if y := point.Y(); y != epoint.Y() {
					t.Errorf("Point: %d: Y: expected %v, got %v\n", idx, epoint.Y(), y)
				}

				if z := point.Z(); z != epoint.Z() {
					t.Errorf("Point: %d: Z: expected %v, got %v\n", idx, epoint.Z(), z)
				}

				if m := point.M(); m != epoint.M() {
					t.Errorf("Point: %d: M: expected %v, got %v\n", idx, epoint.M(), m)
				}
			}
		})
	}
}

func TestLineString_String(t *testing.T) {
	cases := []struct {
		name     string
		line     LineString
		expected string
	}{
		{
			"simple",
			NewLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(1, 3),
					geo.NewPoint(2, 4),
				}),
			),
			"LINESTRING(1 3,2 4)",
		},
		{
			"with Z dimension",
			NewLineString(
				NewBase(NDR, true, false, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointZ(1, 3, 7),
					geo.NewPointZ(2, 4, 5),
				}),
			),
			"LINESTRING(1 3 7,2 4 5)",
		},
		{
			"with M dimension",
			NewLineString(
				NewBase(NDR, false, true, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointM(1, 3, 7),
					geo.NewPointM(2, 4, 5),
				}),
			),
			"LINESTRINGM(1 3 7,2 4 5)",
		},
		{
			"with Z and M dimensions",
			NewLineString(
				NewBase(NDR, true, true, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointZM(1, 3, 7, 2),
					geo.NewPointZM(2, 4, 5, 0),
				}),
			),
			"LINESTRING(1 3 7 2,2 4 5 0)",
		},
		{
			"with SRID",
			NewLineString(
				NewBase(NDR, false, false, true, 4326),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(1, 3),
					geo.NewPoint(2, 4),
				}),
			),
			"SRID=4326;LINESTRING(1 3,2 4)",
		},
		{
			"empty",
			NewLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiPoint([]geo.Point{}),
			),
			"LINESTRING EMPTY",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := c.line.String()

			if s != c.expected {
				t.Errorf("Expected %q, got %q\n", c.expected, s)
			}
		})
	}
}

func TestLineString_Value(t *testing.T) {
	cases := []struct {
		name     string
		line     LineString
		expected string
	}{
		{
			"simple",
			NewLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(1, 3),
					geo.NewPoint(2, 4),
				}),
			),
			"LINESTRING(1 3,2 4)",
		},
		{
			"with Z dimension",
			NewLineString(
				NewBase(NDR, true, false, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointZ(1, 3, 7),
					geo.NewPointZ(2, 4, 5),
				}),
			),
			"LINESTRING(1 3 7,2 4 5)",
		},
		{
			"with M dimension",
			NewLineString(
				NewBase(NDR, false, true, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointM(1, 3, 7),
					geo.NewPointM(2, 4, 5),
				}),
			),
			"LINESTRINGM(1 3 7,2 4 5)",
		},
		{
			"with Z and M dimensions",
			NewLineString(
				NewBase(NDR, true, true, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointZM(1, 3, 7, 2),
					geo.NewPointZM(2, 4, 5, 0),
				}),
			),
			"LINESTRING(1 3 7 2,2 4 5 0)",
		},
		{
			"with SRID",
			NewLineString(
				NewBase(NDR, false, false, true, 4326),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(1, 3),
					geo.NewPoint(2, 4),
				}),
			),
			"SRID=4326;LINESTRING(1 3,2 4)",
		},
		{
			"empty",
			NewLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiPoint([]geo.Point{}),
			),
			"LINESTRING EMPTY",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, err := c.line.Value()
			if err != nil {
				t.Fatalf("Expected not errors, got %v\n", err)
			}

			if s != c.expected {
				t.Errorf("Expected %q, got %q\n", c.expected, s)
			}
		})
	}
}

func TestLineString_MarshalBinary(t *testing.T) {
	cases := []struct {
		name     string
		expected []byte
		data     LineString
	}{
		{
			"simple",
			[]byte{
				1, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0,
				0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64,
			},
			NewLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(1, 3),
					geo.NewPoint(2, 4),
				}),
			),
		},
		{
			"with Z dimension",
			[]byte{
				1, 2, 0, 0, 128, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0,
				0, 28, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0,
				0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 20, 64,
			},
			NewLineString(
				NewBase(NDR, true, false, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointZ(1, 3, 7),
					geo.NewPointZ(2, 4, 5),
				}),
			),
		},
		{
			"with M dimension",
			[]byte{
				1, 2, 0, 0, 64, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0,
				0, 28, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0,
				0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 20, 64,
			},
			NewLineString(
				NewBase(NDR, false, true, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointM(1, 3, 7),
					geo.NewPointM(2, 4, 5),
				}),
			),
		},
		{
			"with Z and M dimension",
			[]byte{
				1, 2, 0, 0, 192, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0,
				0, 0, 28, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0,
				0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0,
				0, 0, 0, 0, 20, 64, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			NewLineString(
				NewBase(NDR, true, true, false, 0),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointZM(1, 3, 7, 2),
					geo.NewPointZM(2, 4, 5, 0),
				}),
			),
		},
		{
			"with SRID",
			[]byte{
				1, 2, 0, 0, 32, 230, 16, 0, 0, 2,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 240, 63, 0,
				0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0,
				0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64,
			},
			NewLineString(
				NewBase(NDR, false, false, true, 4326),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(1, 3),
					geo.NewPoint(2, 4),
				}),
			),
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
