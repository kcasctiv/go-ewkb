package ewkb

import (
	"encoding/binary"

	"github.com/kcasctiv/go-ewkb/geo"
)

func headerSize(hasSRID bool) int {
	if hasSRID {
		return 9
	}

	return 5
}

func pointSize(hasZ, hasM bool) int {
	s := 16
	if hasZ {
		s += 8
	}

	if hasM {
		s += 8
	}

	return s
}

func multiPointSize(mp geo.MultiPoint, hasZ, hasM bool) int {
	return pointSize(hasZ, hasM)*mp.Len() + 4
}

func polygonSize(poly geo.Polygon, hasZ, hasM bool) int {
	size := 4
	for idx := 0; idx < poly.Len(); idx++ {
		size += multiPointSize(poly.Ring(idx), hasZ, hasM)
	}

	return size
}

func multiPolygonSize(mp geo.MultiPolygon, hasZ, hasM bool) int {
	size := 4
	for idx := 0; idx < mp.Len(); idx++ {
		size += polygonSize(mp.Polygon(idx), hasZ, hasM)
	}

	return size
}

func multiLineSize(ml geo.MultiLine, hasZ, hasM bool) int {
	size := 4
	for idx := 0; idx < ml.Len(); idx++ {
		size += multiPointSize(ml.Line(idx), hasZ, hasM)
	}

	return size
}

func collectionSize(geoms []Geometry, hasZ, hasM bool) int {
	size := 4 + headerSize(false)*len(geoms)
	for _, geom := range geoms {
		switch g := geom.(type) {
		case *Point:
			size += pointSize(hasZ, hasM)
		case *LineString:
			size += multiPointSize(g, hasZ, hasM)
		case *MultiPoint:
			size += multiPointSize(g, hasZ, hasM)
		case *Polygon:
			size += polygonSize(g, hasZ, hasM)
		case *MultiLineString:
			size += multiLineSize(g, hasZ, hasM)
		case *MultiPolygon:
			size += multiPolygonSize(g, hasZ, hasM)
		}
	}

	return size
}

func writeHeader(
	h header,
	byteOrder binary.ByteOrder,
	hasSRID bool,
	b []byte,
) int {
	b[0] = h.byteOrder
	byteOrder.PutUint32(b[1:], h.wkbType)

	if !hasSRID {
		return 5
	}

	byteOrder.PutUint32(b[5:], uint32(h.srid))

	return 9
}

func writePoint(
	p geo.Point,
	byteOrder binary.ByteOrder,
	hasZ, hasM bool,
	b []byte,
) int {
	byteOrder.PutUint64(b, uint64(p.X()))
	offset := 8

	byteOrder.PutUint64(b[offset:], uint64(p.Y()))
	offset += 8

	if hasZ {
		byteOrder.PutUint64(b[offset:], uint64(p.Z()))
		offset += 8
	}

	if hasM {
		byteOrder.PutUint64(b[offset:], uint64(p.M()))
		offset += 8
	}

	return offset
}

func writeMultiPoint(
	mp geo.MultiPoint,
	byteOrder binary.ByteOrder,
	hasZ, hasM bool,
	b []byte,
) int {
	byteOrder.PutUint32(b, uint32(mp.Len()))
	offset := 4

	for idx := 0; idx < mp.Len(); idx++ {
		offset += writePoint(mp.Point(idx), byteOrder, hasZ, hasM, b[offset:])
	}

	return offset
}

func writePolygon(
	p geo.Polygon,
	byteOrder binary.ByteOrder,
	hasZ, hasM bool,
	b []byte,
) int {
	byteOrder.PutUint32(b, uint32(p.Len()))
	offset := 4

	for idx := 0; idx < p.Len(); idx++ {
		offset += writeMultiPoint(p.Ring(idx), byteOrder, hasZ, hasM, b[offset:])
	}

	return offset
}

func writeMultiPolygon(
	p geo.MultiPolygon,
	byteOrder binary.ByteOrder,
	hasZ, hasM bool,
	b []byte,
) int {
	byteOrder.PutUint32(b, uint32(p.Len()))
	offset := 4

	for idx := 0; idx < p.Len(); idx++ {
		offset += writePolygon(p.Polygon(idx), byteOrder, hasZ, hasM, b[offset:])
	}

	return offset
}

func writeMultiLine(
	l geo.MultiLine,
	byteOrder binary.ByteOrder,
	hasZ, hasM bool,
	b []byte,
) int {
	byteOrder.PutUint32(b, uint32(l.Len()))
	offset := 4

	for idx := 0; idx < l.Len(); idx++ {
		offset += writeMultiPoint(l.Line(idx), byteOrder, hasZ, hasM, b[offset:])
	}

	return offset
}

func writeCollection(
	geoms []Geometry,
	byteOrder binary.ByteOrder,
	hasZ, hasM bool,
	b []byte,
) int {
	byteOrder.PutUint32(b, uint32(len(geoms)))
	offset := 4

	for _, geom := range geoms {
		switch g := geom.(type) {
		case *Point:
			offset += writePoint(g, byteOrder, hasZ, hasM, b[offset:])
		case *LineString:
			offset += writeMultiPoint(g, byteOrder, hasZ, hasM, b[offset:])
		case *MultiPoint:
			offset += writeMultiPoint(g, byteOrder, hasZ, hasM, b[offset:])
		case *Polygon:
			offset += writePolygon(g, byteOrder, hasZ, hasM, b[offset:])
		case *MultiLineString:
			offset += writeMultiLine(g, byteOrder, hasZ, hasM, b[offset:])
		case *MultiPolygon:
			offset += writeMultiPolygon(g, byteOrder, hasZ, hasM, b[offset:])
		}
	}

	return offset
}
