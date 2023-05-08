package wsdlgen

import (
	_ "embed"
	"github.com/smartlet/wsdl2go/builtin"
	"strings"
)

type Type any // BuiltinType/SimpleType/ComplexType

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

func FieldTypeName(t Type) string {
	switch t := t.(type) {
	case nil:
		return "any"
	case BuiltinType:
		return Identifier(string(t))
	case *SimpleType:
		return Identifier(t.Name)
	case *ComplexType:
		// FIXBUG: any接口
		if t.Base == nil && len(t.Attributes) == 0 && len(t.Elements) == 0 {
			return Identifier(t.Name)
		}
		return "*" + Identifier(t.Name)
	}
	panic("invalid type")
}

func ElementFieldTypeName(e *Element) string {

	var p string

	if e.MaxOccurs == "" || e.MaxOccurs == "0" || e.MaxOccurs == "1" {
		p = ""
	} else {
		p = "[]"
	}
	return p + FieldTypeName(e.Type)
}

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

func Builtin(ns, name string) BuiltinType {
	if strings.ToLower(ns) == "http://www.w3.org/2001/xmlschema" {
		return BuiltinType(builtin.Type(name))
	}
	return ""
}

var defaultPrefix = map[string]string{
	"http://schemas.xmlsoap.org/wsdl/":                             "w",
	"http://schemas.xmlsoap.org/wsdl/soap/":                        "s",
	"http://schemas.microsoft.com/exchange/services/2006/messages": "m",
	"http://schemas.microsoft.com/exchange/services/2006/types":    "t",
	"http://www.w3.org/2001/XMLSchema":                             "x",
}
