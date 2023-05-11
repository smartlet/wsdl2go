package wsdlgen

import (
	_ "embed"
	"github.com/smartlet/wsdl2go/builtin"
	"strings"
)

type Type any // BuiltinType/SimpleType/ComplexType

type BuiltinType string

type SimpleType struct {
	Ns             string // namespace
	Name           string // local name. 为空表示处理过程被删除掉了.
	Base           Type   // BuiltinType/SimpleType
	MinExclusive   []string
	MinInclusive   []string
	MaxExclusive   []string
	MaxInclusive   []string
	TotalDigits    []string
	FractionDigits []string
	Length         []string
	MinLength      []string
	MaxLength      []string
	WhiteSpace     []string
	Pattern        []string
	Enumeration    []string
	List           Type   // list的itemType
	Union          []Type // union的列表
	deprecated     bool   // 表示处理过程被舍弃
}

type ComplexType struct {
	Ns             string // namespace
	Name           string // local name. 为空表示处理过程被删除掉了.
	Base           Type
	MinExclusive   []string
	MinInclusive   []string
	MaxExclusive   []string
	MaxInclusive   []string
	TotalDigits    []string
	FractionDigits []string
	Length         []string
	MinLength      []string
	MaxLength      []string
	WhiteSpace     []string
	Pattern        []string
	Enumeration    []string
	Attributes     []*Attribute
	Elements       []*Element
	deprecated     bool // 表示处理过程被舍弃
}

type Attribute struct {
	Ns      string
	Name    string
	Default string
	Fixed   string
	Use     string
	Type    Type // BuiltinType/SimpleType
}

type Element struct {
	Ns        string
	Name      string
	Default   string
	Fixed     string
	Use       string
	MaxOccurs string
	Type      Type // BuiltinType/SimpleType/ComplexType
}

type Message struct {
	Ns    string
	Name  string
	Parts *NamedSlice[*Element]
}

type Binding struct {
	Ns         string
	Name       string
	PortType   *PortType
	Operations *NamedSlice[*Operation]
}

type Operation struct {
	Ns           string
	Name         string
	Input        *Message
	Output       *Message
	SoapAction11 string
	SoapAction12 string
	InputHeader  []*Element
	OutputHeader []*Element
	InputBody    *Element
	OutputBody   *Element
}

type PortType struct {
	Ns         string
	Name       string
	Operations *NamedSlice[*Operation]
}

func Builtin(ns, name string) BuiltinType {
	if strings.ToLower(ns) == "http://www.w3.org/2001/xmlschema" {
		return BuiltinType(builtin.Type(name))
	}
	return ""
}

func TypeName(t Type) string {
	switch t := t.(type) {
	case nil:
		return "any"
	case BuiltinType:
		return Identifier(string(t))
	case *SimpleType:
		return Identifier(t.Name)
	case *ComplexType:
		return Identifier(t.Name)
	}
	panic("invalid type")
}

func PointerTypeName(t Type) string {
	switch t := t.(type) {
	case nil:
		return "any"
	case BuiltinType:
		return Identifier(string(t))
	case *SimpleType:
		return Identifier(t.Name)
	case *ComplexType:
		if degraded(t) {
			return Identifier(t.Name)
		}
		return "*" + Identifier(t.Name)
	}
	panic("invalid type")
}

func degraded(ct *ComplexType) bool {
	// ComplexType能退化为SimpleType的条件: 无attribute, 无element, base非空且类型为BuiltinType或*SimpleType
	// base为空是abstract=true的.
	if len(ct.Attributes) == 0 && len(ct.Elements) == 0 {
		switch ct.Base.(type) {
		case BuiltinType:
			return true
		case *SimpleType:
			return true
		}
	}
	return false
}
