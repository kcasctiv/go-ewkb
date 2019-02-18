package ewkb

import (
	"bytes"
	"testing"
	"time"

	"github.com/kcasctiv/go-ewkb/geo"
)

func TestNewPolygon(t *testing.T) {
	cases := []struct {
		name string
		base Base
		poly geo.Polygon
	}{
		{
			"little endian",
			NewBase(NDR, false, false, false, 0),
			geo.NewPolygon([]geo.MultiPoint{
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(6, 5),
					geo.NewPoint(7, 8),
					geo.NewPoint(9, 10),
				}),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(4, 3),
					geo.NewPoint(5, 2),
					geo.NewPoint(8, 4),
				}),
			}),
		},
		{
			"big endian",
			NewBase(XDR, false, false, false, 0),
			geo.NewPolygon([]geo.MultiPoint{
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(6, 5),
					geo.NewPoint(7, 8),
					geo.NewPoint(9, 10),
				}),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(4, 3),
					geo.NewPoint(5, 2),
					geo.NewPoint(8, 4),
				}),
			}),
		},
		{
			"has Z",
			NewBase(NDR, true, false, false, 0),
			geo.NewPolygon([]geo.MultiPoint{
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointZ(6, 5, 4),
					geo.NewPointZ(7, 8, 2),
					geo.NewPointZ(9, 10, 1),
				}),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointZ(4, 3, 5),
					geo.NewPointZ(5, 2, 7),
					geo.NewPointZ(8, 4, 8),
				}),
			}),
		},
		{
			"has M",
			NewBase(NDR, false, true, false, 0),
			geo.NewPolygon([]geo.MultiPoint{
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointM(6, 5, 4),
					geo.NewPointM(7, 8, 2),
					geo.NewPointM(9, 10, 1),
				}),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointM(4, 3, 5),
					geo.NewPointM(5, 2, 7),
					geo.NewPointM(8, 4, 8),
				}),
			}),
		},
		{
			"has Z and M",
			NewBase(NDR, true, true, false, 0),
			geo.NewPolygon([]geo.MultiPoint{
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointZM(6, 5, 4, 3),
					geo.NewPointZM(7, 8, 2, 6),
					geo.NewPointZM(9, 10, 1, 2),
				}),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPointZM(4, 3, 5, 7),
					geo.NewPointZM(5, 2, 7, 9),
					geo.NewPointZM(8, 4, 8, 2),
				}),
			}),
		},
		{
			"has SRID",
			NewBase(NDR, false, false, true, 432),
			geo.NewPolygon([]geo.MultiPoint{
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(6, 5),
					geo.NewPoint(7, 8),
					geo.NewPoint(9, 10),
				}),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(4, 3),
					geo.NewPoint(5, 2),
					geo.NewPoint(8, 4),
				}),
			}),
		},
		{
			"has 3 rings",
			NewBase(NDR, false, false, false, 0),
			geo.NewPolygon([]geo.MultiPoint{
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(6, 5),
					geo.NewPoint(7, 8),
					geo.NewPoint(9, 10),
				}),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(4, 3),
					geo.NewPoint(5, 2),
					geo.NewPoint(8, 4),
				}),
				geo.NewMultiPoint([]geo.Point{
					geo.NewPoint(7, 2),
					geo.NewPoint(1, 4),
					geo.NewPoint(9, 3),
				}),
			}),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			poly := NewPolygon(c.base, c.poly)

			checkPolygon(&poly, c.base, c.poly, t)
		})
	}
}

func TestPolygon_UnmarshalBinary(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		valid    bool
		expected Polygon
	}{
		{
			"simple",
			[]byte{
				1, 3, 0, 0, 0, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0, 64,
				0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0,
				0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0,
				0, 0, 8, 64,
			},
			true,
			NewPolygon(
				NewBase(NDR, false, false, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
						geo.NewPoint(1, 3),
					}),
				}),
			),
		},
		{
			"with Z dimension",
			[]byte{
				1, 3, 0, 0, 128, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63,
				0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0,
				0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0,
				0, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 240, 63,
				0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63,
			},
			true,
			NewPolygon(
				NewBase(NDR, true, false, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZ(1, 3, 1),
						geo.NewPointZ(2, 4, 1),
						geo.NewPointZ(3, 2, 1),
						geo.NewPointZ(1, 3, 1),
					}),
				}),
			),
		},
		{
			"with M dimension",
			[]byte{
				1, 3, 0, 0, 64, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240,
				63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0,
				0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0,
				0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0,
				0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63,
			},
			true,
			NewPolygon(
				NewBase(NDR, false, true, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointM(1, 3, 1),
						geo.NewPointM(2, 4, 1),
						geo.NewPointM(3, 2, 1),
						geo.NewPointM(1, 3, 1),
					}),
				}),
			),
		},
		{
			"with Z and M dimension",
			[]byte{
				1, 3, 0, 0, 192, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240,
				63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0,
				0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0,
				16, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0,
				0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0,
				0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 240,
				63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0,
				0, 0, 0, 0, 0, 0, 64,
			},
			true,
			NewPolygon(
				NewBase(NDR, true, true, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZM(1, 3, 1, 2),
						geo.NewPointZM(2, 4, 1, 2),
						geo.NewPointZM(3, 2, 1, 2),
						geo.NewPointZM(1, 3, 1, 2),
					}),
				}),
			),
		},
		{
			"with SRID",
			[]byte{
				1, 3, 0, 0, 32, 230, 16, 0, 0, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0,
				64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0,
				0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8,
				64,
			},
			true,
			NewPolygon(
				NewBase(NDR, false, false, true, 4326),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
						geo.NewPoint(1, 3),
					}),
				}),
			),
		},
		{
			"not a polygon",
			[]byte{1, 2, 0, 0, 0, 0, 0, 0, 0},
			false,
			Polygon{},
		},
		{
			"simple corrupted",
			[]byte{
				1, 3, 0, 0, 0, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0, 64,
				0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0,
				0, 0, 8, 64,
			},
			false,
			Polygon{},
		},
		{
			"with Z dimension corrupted",
			[]byte{
				1, 3, 0, 0, 128, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63,
				0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0,
				0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0,
				0, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 240, 63,
				0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63,
			},
			false,
			Polygon{},
		},
		{
			"with M dimension corrupted",
			[]byte{
				1, 3, 0, 0, 64, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240,
				63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0,
				0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0,
				0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0,
				0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63,
			},
			false,
			Polygon{},
		},
		{
			"with Z and M dimension corrupted",
			[]byte{
				1, 3, 0, 0, 192, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240,
				63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0,
				0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0,
				16, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0,
				0, 0, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0,
				0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 240,
				63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0,
				0, 0, 0, 0, 0, 0, 64,
			},
			false,
			Polygon{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var poly Polygon
			err := poly.UnmarshalBinary(c.data)
			if err != nil && c.valid {
				t.Fatalf("Expected: no errors, got error: %v\n", err)
			}
			if err == nil && !c.valid {
				t.Fatal("Expected: error, got: no errors\n")
			}
			if !c.valid {
				return
			}

			checkPolygon(&poly, &c.expected, &c.expected, t)
		})
	}
}

func TestPolygon_Scan(t *testing.T) {
	cases := []struct {
		name     string
		src      interface{}
		valid    bool
		expected Polygon
	}{
		{
			"binary",
			[]byte{
				1, 3, 0, 0, 0, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0, 64,
				0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0,
				0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0,
				0, 0, 8, 64,
			},
			true,
			NewPolygon(
				NewBase(NDR, false, false, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
						geo.NewPoint(1, 3),
					}),
				}),
			),
		},
		{
			"hex binary",
			[]byte("01030000000100000004000000000000000000F03F00000000000008400000000000000040000000000000104000000000000008400000000000000040000000000000F03F0000000000000840"),
			true,
			NewPolygon(
				NewBase(NDR, false, false, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
						geo.NewPoint(1, 3),
					}),
				}),
			),
		},
		{
			"hex string",
			"01030000000100000004000000000000000000F03F00000000000008400000000000000040000000000000104000000000000008400000000000000040000000000000F03F0000000000000840",
			true,
			NewPolygon(
				NewBase(NDR, false, false, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
						geo.NewPoint(1, 3),
					}),
				}),
			),
		},
		{
			"not valid hex string",
			"01030000000100000004000000000000000000F03G00000000000008400000000000000040000000000000104000000000000008400000000000000040000000000000F03F0000000000000840",
			false,
			Polygon{},
		},
		{
			"not valid hex binary",
			[]byte("01030000000100000004000000000000000000F03G00000000000008400000000000000040000000000000104000000000000008400000000000000040000000000000F03F0000000000000840"),
			false,
			Polygon{},
		},
		{
			"not valid data type",
			time.Now(),
			false,
			Polygon{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var poly Polygon
			err := poly.Scan(c.src)
			if err != nil && c.valid {
				t.Fatalf("Expected: no errors, got error: %v\n", err)
			}
			if err == nil && !c.valid {
				t.Fatal("Expected: error, got: no errors\n")
			}
			if !c.valid {
				return
			}

			checkPolygon(&poly, &c.expected, &c.expected, t)
		})
	}
}

func TestPolygon_String(t *testing.T) {
	cases := []struct {
		name     string
		poly     Polygon
		expected string
	}{
		{
			"simple",
			NewPolygon(
				NewBase(NDR, false, false, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
						geo.NewPoint(1, 3),
					}),
				}),
			),
			"POLYGON((1 3,2 4,3 2,1 3))",
		},
		{
			"with Z dimension",
			NewPolygon(
				NewBase(NDR, true, false, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZ(1, 3, 1),
						geo.NewPointZ(2, 4, 1),
						geo.NewPointZ(3, 2, 1),
						geo.NewPointZ(1, 3, 1),
					}),
				}),
			),
			"POLYGON((1 3 1,2 4 1,3 2 1,1 3 1))",
		},
		{
			"with M dimension",
			NewPolygon(
				NewBase(NDR, false, true, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointM(1, 3, 1),
						geo.NewPointM(2, 4, 1),
						geo.NewPointM(3, 2, 1),
						geo.NewPointM(1, 3, 1),
					}),
				}),
			),
			"POLYGONM((1 3 1,2 4 1,3 2 1,1 3 1))",
		},
		{
			"with Z and M dimensions",
			NewPolygon(
				NewBase(NDR, true, true, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZM(1, 3, 1, 2),
						geo.NewPointZM(2, 4, 1, 2),
						geo.NewPointZM(3, 2, 1, 2),
						geo.NewPointZM(1, 3, 1, 2),
					}),
				}),
			),
			"POLYGON((1 3 1 2,2 4 1 2,3 2 1 2,1 3 1 2))",
		},
		{
			"with SRID",
			NewPolygon(
				NewBase(NDR, false, false, true, 4326),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
						geo.NewPoint(1, 3),
					}),
				}),
			),
			"SRID=4326;POLYGON((1 3,2 4,3 2,1 3))",
		},
		{
			"empty",
			NewPolygon(
				NewBase(NDR, false, false, false, 0),
				geo.NewPolygon([]geo.MultiPoint{}),
			),
			"POLYGON EMPTY",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := c.poly.String()

			if s != c.expected {
				t.Errorf("Expected %q, got %q\n", c.expected, s)
			}
		})
	}
}

func TestPolygon_Value(t *testing.T) {
	cases := []struct {
		name     string
		poly     Polygon
		expected string
	}{
		{
			"simple",
			NewPolygon(
				NewBase(NDR, false, false, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
						geo.NewPoint(1, 3),
					}),
				}),
			),
			"POLYGON((1 3,2 4,3 2,1 3))",
		},
		{
			"with Z dimension",
			NewPolygon(
				NewBase(NDR, true, false, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZ(1, 3, 1),
						geo.NewPointZ(2, 4, 1),
						geo.NewPointZ(3, 2, 1),
						geo.NewPointZ(1, 3, 1),
					}),
				}),
			),
			"POLYGON((1 3 1,2 4 1,3 2 1,1 3 1))",
		},
		{
			"with M dimension",
			NewPolygon(
				NewBase(NDR, false, true, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointM(1, 3, 1),
						geo.NewPointM(2, 4, 1),
						geo.NewPointM(3, 2, 1),
						geo.NewPointM(1, 3, 1),
					}),
				}),
			),
			"POLYGONM((1 3 1,2 4 1,3 2 1,1 3 1))",
		},
		{
			"with Z and M dimensions",
			NewPolygon(
				NewBase(NDR, true, true, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZM(1, 3, 1, 2),
						geo.NewPointZM(2, 4, 1, 2),
						geo.NewPointZM(3, 2, 1, 2),
						geo.NewPointZM(1, 3, 1, 2),
					}),
				}),
			),
			"POLYGON((1 3 1 2,2 4 1 2,3 2 1 2,1 3 1 2))",
		},
		{
			"with SRID",
			NewPolygon(
				NewBase(NDR, false, false, true, 4326),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
						geo.NewPoint(1, 3),
					}),
				}),
			),
			"SRID=4326;POLYGON((1 3,2 4,3 2,1 3))",
		},
		{
			"empty",
			NewPolygon(
				NewBase(NDR, false, false, false, 0),
				geo.NewPolygon([]geo.MultiPoint{}),
			),
			"POLYGON EMPTY",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, err := c.poly.Value()
			if err != nil {
				t.Fatalf("Expected not errors, got %v\n", err)
			}

			if s != c.expected {
				t.Errorf("Expected %q, got %q\n", c.expected, s)
			}
		})
	}
}

func TestPolygon_MarshalBinary(t *testing.T) {
	cases := []struct {
		name     string
		expected []byte
		data     Polygon
	}{
		{
			"simple",
			[]byte{
				1, 3, 0, 0, 0, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0, 64,
				0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0,
				0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0,
				0, 0, 8, 64,
			},
			NewPolygon(
				NewBase(NDR, false, false, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
						geo.NewPoint(1, 3),
					}),
				}),
			),
		},
		{
			"with Z dimension",
			[]byte{
				1, 3, 0, 0, 128, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63,
				0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0,
				0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0,
				0, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 240, 63,
				0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63,
			},
			NewPolygon(
				NewBase(NDR, true, false, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZ(1, 3, 1),
						geo.NewPointZ(2, 4, 1),
						geo.NewPointZ(3, 2, 1),
						geo.NewPointZ(1, 3, 1),
					}),
				}),
			),
		},
		{
			"with M dimension",
			[]byte{
				1, 3, 0, 0, 64, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240,
				63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0,
				0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0,
				0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0,
				0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63,
			},
			NewPolygon(
				NewBase(NDR, false, true, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointM(1, 3, 1),
						geo.NewPointM(2, 4, 1),
						geo.NewPointM(3, 2, 1),
						geo.NewPointM(1, 3, 1),
					}),
				}),
			),
		},
		{
			"with Z and M dimension",
			[]byte{
				1, 3, 0, 0, 192, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240,
				63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0,
				0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0,
				16, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0,
				0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0,
				0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 240,
				63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0,
				0, 0, 0, 0, 0, 0, 64,
			},
			NewPolygon(
				NewBase(NDR, true, true, false, 0),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZM(1, 3, 1, 2),
						geo.NewPointZM(2, 4, 1, 2),
						geo.NewPointZM(3, 2, 1, 2),
						geo.NewPointZM(1, 3, 1, 2),
					}),
				}),
			),
		},
		{
			"with SRID",
			[]byte{
				1, 3, 0, 0, 32, 230, 16, 0, 0, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0,
				64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0,
				0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8,
				64,
			},
			NewPolygon(
				NewBase(NDR, false, false, true, 4326),
				geo.NewPolygon([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
						geo.NewPoint(1, 3),
					}),
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

func checkPolygon(poly *Polygon, ebase Base, epoly geo.Polygon, t *testing.T) {
	if typ := poly.Type(); typ != PolygonType {
		t.Errorf("Type: expected %v, got %v\n", PolygonType, typ)
	}

	if byteOrder := poly.ByteOrder(); byteOrder != ebase.ByteOrder() {
		t.Errorf("ByteOrder: expected %v, got %v\n", ebase.ByteOrder(), byteOrder)
	}

	if hasZ := poly.HasZ(); hasZ != ebase.HasZ() {
		t.Errorf("HasZ: expected %v, got %v\n", ebase.HasZ(), hasZ)
	}

	if hasM := poly.HasM(); hasM != ebase.HasM() {
		t.Errorf("HasM: expected %v, got %v\n", ebase.HasM(), hasM)
	}

	if hasSRID := poly.HasSRID(); hasSRID != ebase.HasSRID() {
		t.Errorf("HasSRID: expected %v, got %v\n", ebase.HasSRID(), hasSRID)
	}

	if srid := poly.SRID(); srid != ebase.SRID() {
		t.Errorf("SRID: expected %v, got %v\n", ebase.SRID(), srid)
	}

	if len := poly.Len(); len != epoly.Len() {
		t.Errorf("Len: expected %v, got %v\n", epoly.Len(), len)
	}

	for idx := 0; idx < poly.Len(); idx++ {
		ring := poly.Ring(idx)
		ering := epoly.Ring(idx)

		if len := ring.Len(); len != ering.Len() {
			t.Errorf("Ring: %d: Len: expected %v, got %v\n", idx, ering.Len(), len)
			continue
		}

		for idx1 := 0; idx1 < ring.Len(); idx1++ {
			point := ring.Point(idx)
			epoint := ering.Point(idx)

			if x := point.X(); x != epoint.X() {
				t.Errorf("Ring: %d: Point: %d: X: expected %v, got %v\n", idx, idx1, epoint.X(), x)
			}

			if y := point.Y(); y != epoint.Y() {
				t.Errorf("Ring: %d: Point: %d: Y: expected %v, got %v\n", idx, idx1, epoint.Y(), y)
			}

			if z := point.Z(); z != epoint.Z() {
				t.Errorf("Ring: %d: Point: %d: Z: expected %v, got %v\n", idx, idx1, epoint.Z(), z)
			}

			if m := point.M(); m != epoint.M() {
				t.Errorf("Ring: %d: Point: %d: M: expected %v, got %v\n", idx, idx1, epoint.M(), m)
			}
		}
	}
}
