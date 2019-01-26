package ewkb

import (
	"database/sql/driver"
	"fmt"
	"math"

	"github.com/kcasctiv/go-ewkb/geo"
)

type MultiLineString struct {
	byteOrder byte
	wkbType   uint32
	srid      int32
	bbox      *bbox
	ml        geo.MultiLine
}

func (l *MultiLineString) ByteOrder() byte             { return l.byteOrder }
func (l *MultiLineString) Type() uint32                { return l.wkbType & uint32(math.MaxUint16) }
func (l *MultiLineString) HasZ() bool                  { return (l.wkbType & zFlag) == zFlag }
func (l *MultiLineString) HasM() bool                  { return (l.wkbType & mFlag) == mFlag }
func (l *MultiLineString) HasSRID() bool               { return (l.wkbType & sridFlag) == sridFlag }
func (l *MultiLineString) HasBBOX() bool               { return (l.wkbType & bboxFlag) == bboxFlag }
func (l *MultiLineString) Line(idx int) geo.MultiPoint { return l.ml.Line(idx) }
func (l *MultiLineString) Len() int                    { return l.ml.Len() }

func (l *MultiLineString) String() string {
	var s string
	if l.HasSRID() {
		s = fmt.Sprintf("SRID=%d;", l.srid)
	}
	s += "MULTILINESTRING "
	if l.HasZ() {
		s += "Z"
	}
	if l.HasM() {
		s += "M"
	}
	s += " ("
	if l.Len() > 0 {
		for idx := 0; idx < l.Len(); idx++ {
			s += printMultiPoint(l.Line(idx), l.HasZ(), l.HasM()) + ", "
		}

		s = s[:len(s)-2]
	}

	return s + ")"
}

func (l *MultiLineString) Scan(src interface{}) error {
	// TODO:
	return nil
}

func (l *MultiLineString) Value() (driver.Value, error) {
	return l.String(), nil
}
