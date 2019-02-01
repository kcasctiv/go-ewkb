package geo

import (
	"fmt"
	"testing"
)

func TestNewPoint(t *testing.T) {
	x := 100.0
	y := -43.32
	p := NewPoint(x, y)

	checkPoint(t, p, x, y, 0, 0, "")
}

func TestNewPointZ(t *testing.T) {
	x := 100.0
	y := -43.32
	z := 66.20
	p := NewPointZ(x, y, z)

	checkPoint(t, p, x, y, z, 0, "")
}

func TestNewPointM(t *testing.T) {
	x := 100.0
	y := -43.32
	m := 66.20
	p := NewPointM(x, y, m)

	checkPoint(t, p, x, y, 0, m, "")
}

func TestNewPointZM(t *testing.T) {
	x := 100.0
	y := -43.32
	z := 12.5
	m := -66.20
	p := NewPointZM(x, y, z, m)

	checkPoint(t, p, x, y, z, m, "")
}

func TestMultiPoint(t *testing.T) {
	points := []Point{
		NewPoint(1, 2),
		NewPointZ(3, 4, 5),
		NewPointM(6, 7, 8),
		NewPointZM(9, 10, 11, 12),
	}

	mp := NewMultiPoint(points)

	if mpl := mp.Len(); mpl != len(points) {
		t.Fatalf("Len: extected %v, got %v\n", len(points), mpl)
	}

	for idx := 0; idx < mp.Len(); idx++ {
		p := points[idx]
		mpp := mp.Point(idx)

		checkPoint(t, mpp, p.X(), p.Y(), p.Z(), p.M(), fmt.Sprintf("Point: %d: ", idx))
	}
}

func TestPolygon(t *testing.T) {
	rings := []MultiPoint{
		NewMultiPoint([]Point{
			NewPoint(1, 2),
			NewPoint(3, 4),
			NewPoint(11, 12),
		}),
		NewMultiPoint([]Point{
			NewPoint(5, 6),
			NewPoint(7, 8),
			NewPoint(71, 81),
			NewPoint(9, 10),
		}),
	}

	p := NewPolygon(rings)

	if pl := p.Len(); pl != len(rings) {
		t.Fatalf("Len: extected %v, got %v\n", len(rings), pl)
	}

	for idx := 0; idx < p.Len(); idx++ {
		ering := rings[idx]
		ring := p.Ring(idx)

		checkMultiPoint(t, ring, ering, fmt.Sprintf("Ring %d: ", idx))
	}
}

func TestNewMultiLine(t *testing.T) {
	lines := []MultiPoint{
		NewMultiPoint([]Point{
			NewPoint(14, 2.2),
			NewPoint(1.3, 2.4),
		}),
		NewMultiPoint([]Point{
			NewPoint(5.4, 6.3),
			NewPoint(7.1, 8.9),
			NewPoint(9.3, 1.01),
		}),
	}

	ml := NewMultiLine(lines)

	if mll := ml.Len(); mll != len(lines) {
		t.Fatalf("Len: extected %v, got %v\n", len(lines), mll)
	}

	for idx := 0; idx < ml.Len(); idx++ {
		eline := lines[idx]
		line := ml.Line(idx)

		checkMultiPoint(t, line, eline, fmt.Sprintf("Line %d: ", idx))
	}
}

func TestNewMultiPolygon(t *testing.T) {
	pols := []Polygon{
		NewPolygon([]MultiPoint{
			NewMultiPoint([]Point{
				NewPoint(5.4, 6.3),
				NewPoint(7.1, 8.9),
				NewPoint(9.3, 1.01),
			}),
		}),
		NewPolygon([]MultiPoint{
			NewMultiPoint([]Point{
				NewPoint(5, 6),
				NewPoint(7, 8),
				NewPoint(71, 81),
				NewPoint(9, 10),
			}),
		}),
	}

	mp := NewMultiPolygon(pols)
	if mpl := mp.Len(); mpl != len(pols) {
		t.Fatalf("Len: extected %v, got %v\n", len(pols), mpl)
	}

	for idx := 0; idx < mp.Len(); idx++ {
		epoly := pols[idx]
		poly := mp.Polygon(idx)

		if pl := poly.Len(); pl != epoly.Len() {
			t.Errorf("Polygon: %d: Len: extected %v, got %v\n", idx, epoly.Len(), pl)
			continue
		}

		for idx1 := 0; idx1 < poly.Len(); idx1++ {
			ering := epoly.Ring(idx1)
			ring := poly.Ring(idx1)

			checkMultiPoint(t, ring, ering, fmt.Sprintf("Polygon %d: Ring %d: ", idx, idx1))
		}
	}
}

func checkPoint(t *testing.T, p Point, x, y, z, m float64, prefix string) {
	if px := p.X(); px != x {
		t.Errorf("%sX: expected %v, got %v\n", prefix, x, px)
	}

	if py := p.Y(); py != y {
		t.Errorf("%sY: expected %v, got %v\n", prefix, y, py)
	}

	if pz := p.Z(); pz != z {
		t.Errorf("%sZ: expected %v, got %v\n", prefix, z, pz)
	}

	if pm := p.M(); pm != m {
		t.Errorf("%sM: expected %v, got %v\n", prefix, m, pm)
	}
}

func checkMultiPoint(t *testing.T, mp, emp MultiPoint, prefix string) {
	if mpl := mp.Len(); mpl != emp.Len() {
		t.Errorf("%sLen: extected %v, got %v\n", prefix, emp.Len(), mpl)
		return
	}

	for idx := 0; idx < mp.Len(); idx++ {
		ep := emp.Point(idx)
		p := mp.Point(idx)

		checkPoint(t, p, ep.X(), ep.Y(), ep.Z(), ep.M(), fmt.Sprintf("%sPoint %d: ", prefix, idx))
	}
}
