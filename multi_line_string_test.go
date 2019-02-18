package ewkb

import (
	"bytes"
	"testing"
	"time"

	"github.com/kcasctiv/go-ewkb/geo"
)

func TestNewMultiLineString(t *testing.T) {
	cases := []struct {
		name  string
		base  Base
		mline geo.MultiLine
	}{
		{
			"little endian",
			NewBase(NDR, false, false, false, 0),
			geo.NewMultiLine([]geo.MultiPoint{
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
			geo.NewMultiLine([]geo.MultiPoint{
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
			geo.NewMultiLine([]geo.MultiPoint{
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
			geo.NewMultiLine([]geo.MultiPoint{
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
			geo.NewMultiLine([]geo.MultiPoint{
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
			geo.NewMultiLine([]geo.MultiPoint{
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
			"has 3 lines",
			NewBase(NDR, false, false, false, 0),
			geo.NewMultiLine([]geo.MultiPoint{
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
			mline := NewMultiLineString(c.base, c.mline)

			checkMultiLine(&mline, c.base, c.mline, t)
		})
	}
}

func TestMultiLineString_UnmarshalBinary(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		valid    bool
		expected MultiLineString
	}{
		{
			"simple",
			[]byte{
				1, 5, 0, 0, 0, 1, 0, 0, 0, 1, 2, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0,
				64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0,
				0, 0, 0, 0, 64,
			},
			true,
			NewMultiLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
					}),
				}),
			),
		},
		{
			"with Z dimension",
			[]byte{
				1, 5, 0, 0, 128, 1, 0, 0, 0, 1, 2, 0, 0, 128, 3, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0,
				0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0,
				0, 0, 64, 0, 0, 0, 0, 0, 0, 240, 63,
			},
			true,
			NewMultiLineString(
				NewBase(NDR, true, false, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZ(1, 3, 1),
						geo.NewPointZ(2, 4, 1),
						geo.NewPointZ(3, 2, 1),
					}),
				}),
			),
		},
		{
			"with M dimension",
			[]byte{
				1, 5, 0, 0, 64, 1, 0, 0, 0, 1, 2, 0, 0, 64, 3, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0,
				0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0,
				0, 0, 64, 0, 0, 0, 0, 0, 0, 240, 63,
			},
			true,
			NewMultiLineString(
				NewBase(NDR, false, true, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointM(1, 3, 1),
						geo.NewPointM(2, 4, 1),
						geo.NewPointM(3, 2, 1),
					}),
				}),
			),
		},
		{
			"with Z and M dimension",
			[]byte{
				1, 5, 0, 0, 192, 1, 0, 0, 0, 1, 2, 0, 0, 192, 3, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0,
				0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0,
				0, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0,
				0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64,
			},
			true,
			NewMultiLineString(
				NewBase(NDR, true, true, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZM(1, 3, 1, 2),
						geo.NewPointZM(2, 4, 1, 2),
						geo.NewPointZM(3, 2, 1, 2),
					}),
				}),
			),
		},
		{
			"with SRID",
			[]byte{
				1, 5, 0, 0, 32, 230, 16, 0, 0, 1, 0, 0, 0, 1, 2, 0, 0, 0, 3, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0,
				0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0,
				8, 64, 0, 0, 0, 0, 0, 0, 0, 64,
			},
			true,
			NewMultiLineString(
				NewBase(NDR, false, false, true, 4326),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
					}),
				}),
			),
		},
		{
			"not a multi line string",
			[]byte{1, 2, 0, 0, 0, 0, 0, 0, 0},
			false,
			MultiLineString{},
		},
		{
			"simple corrupted",
			[]byte{
				1, 5, 0, 0, 0, 1, 0, 0, 0, 1, 2, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0,
				64, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0,
				0, 0, 0, 0, 64,
			},
			false,
			MultiLineString{},
		},
		{
			"with Z dimension corrupted",
			[]byte{
				1, 5, 0, 0, 128, 1, 0, 0, 0, 1, 2, 0, 0, 128, 3, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 64, 0,
				0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0,
				0, 0, 64, 0, 0, 0, 0, 0, 0, 240, 63,
			},
			false,
			MultiLineString{},
		},
		{
			"with M dimension corrupted",
			[]byte{
				1, 5, 0, 0, 64, 1, 0, 0, 0, 1, 2, 0, 0, 64, 3, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0,
				0, 0, 0, 0, 0, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0,
				0, 0, 64, 0, 0, 0, 0, 0, 0, 240, 63,
			},
			false,
			MultiLineString{},
		},
		{
			"with Z and M dimension corrupted",
			[]byte{
				1, 5, 0, 0, 192, 1, 0, 0, 0, 1, 2, 0, 0, 192, 3, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0,
				0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0,
				0, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0,
				0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64,
			},
			false,
			MultiLineString{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var mline MultiLineString
			err := mline.UnmarshalBinary(c.data)
			if err != nil && c.valid {
				t.Fatalf("Expected: no errors, got error: %v\n", err)
			}
			if err == nil && !c.valid {
				t.Fatal("Expected: error, got: no errors\n")
			}
			if !c.valid {
				return
			}

			checkMultiLine(&mline, &c.expected, &c.expected, t)
		})
	}
}

func TestMultiLineString_Scan(t *testing.T) {
	cases := []struct {
		name     string
		src      interface{}
		valid    bool
		expected MultiLineString
	}{
		{
			"binary",
			[]byte{
				1, 5, 0, 0, 0, 1, 0, 0, 0, 1, 2, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0,
				64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0,
				0, 0, 0, 0, 64,
			},
			true,
			NewMultiLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
					}),
				}),
			),
		},
		{
			"hex binary",
			[]byte("010500000001000000010200000003000000000000000000F03F00000000000008400000000000000040000000000000104000000000000008400000000000000040"),
			true,
			NewMultiLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
					}),
				}),
			),
		},
		{
			"hex string",
			"010500000001000000010200000003000000000000000000F03F00000000000008400000000000000040000000000000104000000000000008400000000000000040",
			true,
			NewMultiLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
					}),
				}),
			),
		},
		{
			"not valid hex string",
			"010500000001000000010200000003000000000000000000H03F00000000000008400000000000000040000000000000104000000000000008400000000000000040",
			false,
			MultiLineString{},
		},
		{
			"not valid hex binary",
			[]byte("010500000001000000010200000003000000000000000000H03F00000000000008400000000000000040000000000000104000000000000008400000000000000040"),
			false,
			MultiLineString{},
		},
		{
			"not valid data type",
			time.Now(),
			false,
			MultiLineString{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var mline MultiLineString
			err := mline.Scan(c.src)
			if err != nil && c.valid {
				t.Fatalf("Expected: no errors, got error: %v\n", err)
			}
			if err == nil && !c.valid {
				t.Fatal("Expected: error, got: no errors\n")
			}
			if !c.valid {
				return
			}

			checkMultiLine(&mline, &c.expected, &c.expected, t)
		})
	}
}

func TestMultiLineString_String(t *testing.T) {
	cases := []struct {
		name     string
		mline    MultiLineString
		expected string
	}{
		{
			"simple",
			NewMultiLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
					}),
				}),
			),
			"MULTILINESTRING((1 3,2 4,3 2))",
		},
		{
			"with Z dimension",
			NewMultiLineString(
				NewBase(NDR, true, false, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZ(1, 3, 1),
						geo.NewPointZ(2, 4, 1),
						geo.NewPointZ(3, 2, 1),
					}),
				}),
			),
			"MULTILINESTRING((1 3 1,2 4 1,3 2 1))",
		},
		{
			"with M dimension",
			NewMultiLineString(
				NewBase(NDR, false, true, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointM(1, 3, 1),
						geo.NewPointM(2, 4, 1),
						geo.NewPointM(3, 2, 1),
					}),
				}),
			),
			"MULTILINESTRINGM((1 3 1,2 4 1,3 2 1))",
		},
		{
			"with Z and M dimensions",
			NewMultiLineString(
				NewBase(NDR, true, true, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZM(1, 3, 1, 2),
						geo.NewPointZM(2, 4, 1, 2),
						geo.NewPointZM(3, 2, 1, 2),
					}),
				}),
			),
			"MULTILINESTRING((1 3 1 2,2 4 1 2,3 2 1 2))",
		},
		{
			"with SRID",
			NewMultiLineString(
				NewBase(NDR, false, false, true, 4326),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
					}),
				}),
			),
			"SRID=4326;MULTILINESTRING((1 3,2 4,3 2))",
		},
		{
			"empty",
			NewMultiLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{}),
			),
			"MULTILINESTRING EMPTY",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := c.mline.String()

			if s != c.expected {
				t.Errorf("Expected %q, got %q\n", c.expected, s)
			}
		})
	}
}

func TestMultiLineString_Value(t *testing.T) {
	cases := []struct {
		name     string
		mline    MultiLineString
		expected string
	}{
		{
			"simple",
			NewMultiLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
					}),
				}),
			),
			"MULTILINESTRING((1 3,2 4,3 2))",
		},
		{
			"with Z dimension",
			NewMultiLineString(
				NewBase(NDR, true, false, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZ(1, 3, 1),
						geo.NewPointZ(2, 4, 1),
						geo.NewPointZ(3, 2, 1),
					}),
				}),
			),
			"MULTILINESTRING((1 3 1,2 4 1,3 2 1))",
		},
		{
			"with M dimension",
			NewMultiLineString(
				NewBase(NDR, false, true, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointM(1, 3, 1),
						geo.NewPointM(2, 4, 1),
						geo.NewPointM(3, 2, 1),
					}),
				}),
			),
			"MULTILINESTRINGM((1 3 1,2 4 1,3 2 1))",
		},
		{
			"with Z and M dimensions",
			NewMultiLineString(
				NewBase(NDR, true, true, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZM(1, 3, 1, 2),
						geo.NewPointZM(2, 4, 1, 2),
						geo.NewPointZM(3, 2, 1, 2),
					}),
				}),
			),
			"MULTILINESTRING((1 3 1 2,2 4 1 2,3 2 1 2))",
		},
		{
			"with SRID",
			NewMultiLineString(
				NewBase(NDR, false, false, true, 4326),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
					}),
				}),
			),
			"SRID=4326;MULTILINESTRING((1 3,2 4,3 2))",
		},
		{
			"empty",
			NewMultiLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{}),
			),
			"MULTILINESTRING EMPTY",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, _ := c.mline.Value()

			if s != c.expected {
				t.Errorf("Expected %q, got %q\n", c.expected, s)
			}
		})
	}
}

func TestMultiLineString_MarshalBinary(t *testing.T) {
	cases := []struct {
		name     string
		expected []byte
		data     MultiLineString
	}{
		{
			"simple",
			[]byte{
				1, 5, 0, 0, 0, 1, 0, 0, 0, 1, 2, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0,
				64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0,
				0, 0, 0, 0, 64,
			},
			NewMultiLineString(
				NewBase(NDR, false, false, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
					}),
				}),
			),
		},
		{
			"with Z dimension",
			[]byte{
				1, 5, 0, 0, 128, 1, 0, 0, 0, 1, 2, 0, 0, 128, 3, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0,
				0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0,
				0, 0, 64, 0, 0, 0, 0, 0, 0, 240, 63,
			},
			NewMultiLineString(
				NewBase(NDR, true, false, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZ(1, 3, 1),
						geo.NewPointZ(2, 4, 1),
						geo.NewPointZ(3, 2, 1),
					}),
				}),
			),
		},
		{
			"with M dimension",
			[]byte{
				1, 5, 0, 0, 64, 1, 0, 0, 0, 1, 2, 0, 0, 64, 3, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0,
				0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0,
				0, 0, 64, 0, 0, 0, 0, 0, 0, 240, 63,
			},
			NewMultiLineString(
				NewBase(NDR, false, true, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointM(1, 3, 1),
						geo.NewPointM(2, 4, 1),
						geo.NewPointM(3, 2, 1),
					}),
				}),
			),
		},
		{
			"with Z and M dimension",
			[]byte{
				1, 5, 0, 0, 192, 1, 0, 0, 0, 1, 2, 0, 0, 192, 3, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0,
				240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0,
				0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0,
				0, 64, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0,
				0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64,
			},
			NewMultiLineString(
				NewBase(NDR, true, true, false, 0),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPointZM(1, 3, 1, 2),
						geo.NewPointZM(2, 4, 1, 2),
						geo.NewPointZM(3, 2, 1, 2),
					}),
				}),
			),
		},
		{
			"with SRID",
			[]byte{
				1, 5, 0, 0, 32, 230, 16, 0, 0, 1, 0, 0, 0, 1, 2, 0, 0, 0, 3, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 8, 64, 0, 0,
				0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 16, 64, 0, 0, 0, 0, 0, 0,
				8, 64, 0, 0, 0, 0, 0, 0, 0, 64,
			},
			NewMultiLineString(
				NewBase(NDR, false, false, true, 4326),
				geo.NewMultiLine([]geo.MultiPoint{
					geo.NewMultiPoint([]geo.Point{
						geo.NewPoint(1, 3),
						geo.NewPoint(2, 4),
						geo.NewPoint(3, 2),
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

func checkMultiLine(
	mline *MultiLineString, ebase Base,
	emline geo.MultiLine, t *testing.T,
) {
	if typ := mline.Type(); typ != MultiLineType {
		t.Errorf("Type: expected %v, got %v\n", MultiLineType, typ)
	}

	if byteOrder := mline.ByteOrder(); byteOrder != ebase.ByteOrder() {
		t.Errorf("ByteOrder: expected %v, got %v\n", ebase.ByteOrder(), byteOrder)
	}

	if hasZ := mline.HasZ(); hasZ != ebase.HasZ() {
		t.Errorf("HasZ: expected %v, got %v\n", ebase.HasZ(), hasZ)
	}

	if hasM := mline.HasM(); hasM != ebase.HasM() {
		t.Errorf("HasM: expected %v, got %v\n", ebase.HasM(), hasM)
	}

	if hasSRID := mline.HasSRID(); hasSRID != ebase.HasSRID() {
		t.Errorf("HasSRID: expected %v, got %v\n", ebase.HasSRID(), hasSRID)
	}

	if srid := mline.SRID(); srid != ebase.SRID() {
		t.Errorf("SRID: expected %v, got %v\n", ebase.SRID(), srid)
	}

	if len := mline.Len(); len != emline.Len() {
		t.Errorf("Len: expected %v, got %v\n", emline.Len(), len)
	}

	for idx := 0; idx < mline.Len(); idx++ {
		line := mline.Line(idx)
		eline := emline.Line(idx)

		if len := line.Len(); len != eline.Len() {
			t.Errorf("Line: %d: Len: expected %v, got %v\n", idx, eline.Len(), len)
			continue
		}

		for idx1 := 0; idx1 < line.Len(); idx1++ {
			point := line.Point(idx)
			epoint := eline.Point(idx)

			if x := point.X(); x != epoint.X() {
				t.Errorf("Line: %d: Point: %d: X: expected %v, got %v\n", idx, idx1, epoint.X(), x)
			}

			if y := point.Y(); y != epoint.Y() {
				t.Errorf("Line: %d: Point: %d: Y: expected %v, got %v\n", idx, idx1, epoint.Y(), y)
			}

			if z := point.Z(); z != epoint.Z() {
				t.Errorf("Line: %d: Point: %d: Z: expected %v, got %v\n", idx, idx1, epoint.Z(), z)
			}

			if m := point.M(); m != epoint.M() {
				t.Errorf("Line: %d: Point: %d: M: expected %v, got %v\n", idx, idx1, epoint.M(), m)
			}
		}
	}
}
