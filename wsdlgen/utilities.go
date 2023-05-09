package wsdlgen

import (
	"github.com/smartlet/wsdl2go/wsdl"
	"regexp"
	"strings"
)

const diff = 'A' - 'a'

func Identifier(v string) string {
	var sb strings.Builder
	sb.Grow(len(v))
	for i, c := range v {
		if 'a' <= c && c <= 'z' {
			// 首字母大写
			if i == 0 {
				sb.WriteRune(c + diff)
			}
			sb.WriteRune(c)
		} else if 'A' <= c && c <= 'Z' {
			sb.WriteRune(c)
		} else if '0' <= c && c <= '9' {
			sb.WriteRune(c)
		} else if c == '_' {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

const defaultCapacity = 8

type NamedSlice[V any] struct {
	slice []V
	names map[string]V
}

func NewNamedSlice[V any]() *NamedSlice[V] {
	return &NamedSlice[V]{
		slice: make([]V, 0, defaultCapacity),
		names: make(map[string]V, defaultCapacity),
	}
}

func (s *NamedSlice[V]) Set(k string, v V) bool {
	if _, ok := s.names[k]; !ok {
		s.slice = append(s.slice, v)
		s.names[k] = v
		return false
	}
	return true
}

func (s *NamedSlice[V]) Get(k string) V {
	v, _ := s.names[k]
	return v
}

func (s *NamedSlice[V]) All() []V {
	return s.slice
}

type NsOrderSlice[V any] struct {
	data map[string][]V
	hash map[string][]V
}

func NewNsOrderSlice[V any]() *NsOrderSlice[V] {
	return &NsOrderSlice[V]{
		data: make(map[string][]V, defaultCapacity),
		hash: make(map[string][]V, defaultCapacity),
	}
}

func (m *NsOrderSlice[V]) Add(ns, k string, v V) {
	set := m.data[ns]
	set = append(set, v)
	m.data[ns] = set

	bag := m.hash[k]
	bag = append(bag, v)
	m.hash[k] = bag
}

func (m *NsOrderSlice[V]) AllByNs() map[string][]V {
	return m.data
}

func (m *NsOrderSlice[V]) AllByKey() map[string][]V {
	return m.hash
}

type NsNamedSlice[V any] struct {
	data map[string]*NamedSlice[V]
	null V
}

func NewNsNamedSlice[V any]() *NsNamedSlice[V] {
	return &NsNamedSlice[V]{
		data: make(map[string]*NamedSlice[V], defaultCapacity),
	}
}

func (m *NsNamedSlice[V]) Set(ns, k string, v V) bool {
	set, ok := m.data[ns]
	if !ok {
		set = NewNamedSlice[V]()
		m.data[ns] = set
	}
	if _, ok = set.names[k]; !ok {
		set.slice = append(set.slice, v)
		set.names[k] = v
		return false
	}
	return true
}
func (m *NsNamedSlice[V]) Get(ns, k string) V {
	if set, ok := m.data[ns]; ok {
		return set.Get(k)
	}
	return m.null
}

func (m *NsNamedSlice[V]) All() map[string]*NamedSlice[V] {
	return m.data
}

func RestrictionSimpleType(b *wsdl.Restriction, st *SimpleType) {
	st.MinExclusive = wsdl.MinExclusiveValues(b.MinExclusive)
	st.MinInclusive = wsdl.MinInclusiveValues(b.MinInclusive)
	st.MaxExclusive = wsdl.MaxExclusiveValues(b.MaxExclusive)
	st.MaxInclusive = wsdl.MaxInclusiveValues(b.MaxInclusive)
	st.TotalDigits = wsdl.TotalDigitsValues(b.TotalDigits)
	st.FractionDigits = wsdl.FractionDigitsValues(b.FractionDigits)
	st.Length = wsdl.LengthValues(b.Length)
	st.MinLength = wsdl.MinLengthValues(b.MinLength)
	st.MaxLength = wsdl.MaxLengthValues(b.MaxLength)
	st.WhiteSpace = wsdl.WhiteSpaceValues(b.WhiteSpace)
	st.Pattern = wsdl.PatternValues(b.Pattern)
	st.Enumeration = wsdl.EnumerationValues(b.Enumeration)
}

func RestrictionComplexType(b *wsdl.Restriction, ct *ComplexType) {
	ct.MinExclusive = wsdl.MinExclusiveValues(b.MinExclusive)
	ct.MinInclusive = wsdl.MinInclusiveValues(b.MinInclusive)
	ct.MaxExclusive = wsdl.MaxExclusiveValues(b.MaxExclusive)
	ct.MaxInclusive = wsdl.MaxInclusiveValues(b.MaxInclusive)
	ct.TotalDigits = wsdl.TotalDigitsValues(b.TotalDigits)
	ct.FractionDigits = wsdl.FractionDigitsValues(b.FractionDigits)
	ct.Length = wsdl.LengthValues(b.Length)
	ct.MinLength = wsdl.MinLengthValues(b.MinLength)
	ct.MaxLength = wsdl.MaxLengthValues(b.MaxLength)
	ct.WhiteSpace = wsdl.WhiteSpaceValues(b.WhiteSpace)
	ct.Pattern = wsdl.PatternValues(b.Pattern)
	ct.Enumeration = wsdl.EnumerationValues(b.Enumeration)
}

var whitespace = regexp.MustCompile(`\s+`)

func nvl(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

func If[V any](c bool, v1 V, v2 V) V {
	if c {
		return v1
	}
	return v2
}
