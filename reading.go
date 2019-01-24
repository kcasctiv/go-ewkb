package ewkb

import (
	"encoding/binary"
	"errors"
	"math"

	"github.com/kcasctiv/go-ewkb/geo"
)

type readPointFunc func(b []byte, byteOrder binary.ByteOrder) (geo.Point, int, error)

func getReadPointFunc(wkbType uint32) readPointFunc {
	if (wkbType & zFlag) == zFlag {
		if (wkbType & mFlag) == mFlag {
			return readPointZM
		}
		return readPointZ
	}
	if (wkbType & mFlag) == mFlag {
		return readPointM
	}

	return readPoint
}

func readPoint(b []byte, byteOrder binary.ByteOrder) (geo.Point, int, error) {
	if len(b) < 16 {
		return nil, 0, errors.New("out of range")
	}

	x := math.Float64frombits(byteOrder.Uint64(b))
	y := math.Float64frombits(byteOrder.Uint64(b[8:]))
	return geo.NewPoint(x, y), 16, nil
}

func readPointZ(b []byte, byteOrder binary.ByteOrder) (geo.Point, int, error) {
	if len(b) < 24 {
		return nil, 0, errors.New("out of range")
	}

	x := math.Float64frombits(byteOrder.Uint64(b))
	y := math.Float64frombits(byteOrder.Uint64(b[8:]))
	z := math.Float64frombits(byteOrder.Uint64(b[16:]))
	return geo.NewPointZ(x, y, z), 24, nil
}

func readPointM(b []byte, byteOrder binary.ByteOrder) (geo.Point, int, error) {
	if len(b) < 24 {
		return nil, 0, errors.New("out of range")
	}

	x := math.Float64frombits(byteOrder.Uint64(b))
	y := math.Float64frombits(byteOrder.Uint64(b[8:]))
	m := math.Float64frombits(byteOrder.Uint64(b[16:]))
	return geo.NewPointM(x, y, m), 24, nil
}

func readPointZM(b []byte, byteOrder binary.ByteOrder) (geo.Point, int, error) {
	if len(b) < 32 {
		return nil, 0, errors.New("out of range")
	}

	x := math.Float64frombits(byteOrder.Uint64(b))
	y := math.Float64frombits(byteOrder.Uint64(b[8:]))
	z := math.Float64frombits(byteOrder.Uint64(b[16:]))
	m := math.Float64frombits(byteOrder.Uint64(b[24:]))
	return geo.NewPointZM(x, y, z, m), 32, nil
}

func readMultiPoint(
	b []byte, byteOrder binary.ByteOrder,
	readFunc readPointFunc) (geo.MultiPoint, int, error) {
	if len(b) < 4 {
		return nil, 0, errors.New("out of range")
	}

	mplen := int(byteOrder.Uint32(b))
	points := make([]geo.Point, mplen)
	var err error
	var n int
	offset := 4
	for idx := 0; idx < mplen; idx++ {
		points[idx], n, err = readFunc(b[offset:], byteOrder)
		if err != nil {
			return nil, 0, err
		}

		offset += n
	}

	return geo.NewMultiPoint(points), offset, nil
}

func readPolygon(b []byte, byteOrder binary.ByteOrder,
	readFunc readPointFunc) (geo.Polygon, int, error) {
	if len(b) < 4 {
		return nil, 0, errors.New("out of range")
	}

	plen := int(byteOrder.Uint32(b))
	rings := make([]geo.MultiPoint, plen)
	var err error
	var n int
	offset := 4
	for idx := 0; idx < plen; idx++ {
		rings[idx], n, err = readMultiPoint(b[offset:], byteOrder, readFunc)
		if err != nil {
			return nil, 0, err
		}
		offset += n
	}

	return geo.NewPolygon(rings), offset, nil
}

func readMultiLine(b []byte, byteOrder binary.ByteOrder,
	readFunc readPointFunc) (geo.MultiLine, int, error) {
	if len(b) < 4 {
		return nil, 0, errors.New("out of range")
	}

	plen := int(byteOrder.Uint32(b))
	lines := make([]geo.MultiPoint, plen)
	var err error
	var n int
	offset := 4
	for idx := 0; idx < plen; idx++ {
		lines[idx], n, err = readMultiPoint(b[offset:], byteOrder, readFunc)
		if err != nil {
			return nil, 0, err
		}
		offset += n
	}

	return geo.NewMultiLine(lines), offset, nil
}

func readMultiPolygon(b []byte, byteOrder binary.ByteOrder,
	readFunc readPointFunc) (geo.MultiPolygon, int, error) {
	if len(b) < 4 {
		return nil, 0, errors.New("out of range")
	}

	mplen := int(byteOrder.Uint32(b))
	pols := make([]geo.Polygon, mplen)
	var err error
	var n int
	offset := 4
	for idx := 0; idx < mplen; idx++ {
		pols[idx], n, err = readPolygon(b[offset:], byteOrder, readFunc)
		if err != nil {
			return nil, 0, err
		}

		offset += n
	}

	return geo.NewMultiPolygon(pols), offset, nil
}
