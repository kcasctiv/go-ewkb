package ewkb

import (
	"database/sql/driver"
	"fmt"

	"github.com/kcasctiv/go-ewkb/geo"
)

type MultiLineString struct {
	header
	ml geo.MultiLine
}

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
