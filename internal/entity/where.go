package entity

import (
	"fmt"
	"strings"
)

type Conjunction int

const (
	ConjAnd Conjunction = iota
	ConjOr
)

type Where struct {
	Conj        Conjunction
	FieldNames  []string
	FieldValues []any
}

func (w *Where) Parse() (string, []any, error) {
	if w == nil {
		return "", nil, fmt.Errorf("invalid where: nil where")
	}
	if len(w.FieldNames) != len(w.FieldValues) {
		return "", nil, fmt.Errorf("invalid where")
	}
	if len(w.FieldNames) == 0 {
		return "", nil, fmt.Errorf("invalid where: empty fields")
	}

	args := make([]any, 0, len(w.FieldValues))
	sb := strings.Builder{}
	sb.WriteString(" WHERE ")
	conj := " AND "
	if w.Conj == ConjOr {
		conj = " OR "
	}
	for i := 0; i < len(w.FieldValues); i++ {
		sb.WriteString(w.FieldNames[i] + " = ?")
		args = append(args, w.FieldValues[i])
		if i < len(w.FieldValues)-1 {
			// do not need to append conj to the tail
			sb.WriteString(conj)
		}
	}

	return sb.String(), args, nil
}
