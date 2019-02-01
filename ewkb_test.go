package ewkb

import "testing"

func TestNewBase(t *testing.T) {
	cases := []struct {
		name                string
		byteOrder           byte
		hasZ, hasM, hasSRID bool
		srid                int32
	}{
		{name: "little endian", byteOrder: NDR},
		{name: "big endian", byteOrder: XDR},
		{name: "has Z", hasZ: true},
		{name: "has M", hasM: true},
		{name: "has SRID", hasSRID: true, srid: 132},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := NewBase(c.byteOrder, c.hasZ, c.hasM, c.hasSRID, c.srid)
			if byteOrder := b.ByteOrder(); byteOrder != c.byteOrder {
				t.Errorf("ByteOrder: expected %v, got %v\n", c.byteOrder, byteOrder)
			}

			if hasZ := b.HasZ(); hasZ != c.hasZ {
				t.Errorf("HasZ: expected %v, got %v\n", c.hasZ, hasZ)
			}

			if hasM := b.HasM(); hasM != c.hasM {
				t.Errorf("HasM: expected %v, got %v\n", c.hasM, hasM)
			}

			if hasSRID := b.HasSRID(); hasSRID != c.hasSRID {
				t.Errorf("HasSRID: expected %v, got %v\n", c.hasSRID, hasSRID)
			}

			if srid := b.SRID(); srid != c.srid {
				t.Errorf("SRID: expected %v, got %v\n", c.srid, srid)
			}
		})
	}
}
