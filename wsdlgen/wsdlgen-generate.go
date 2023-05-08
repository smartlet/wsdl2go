package wsdlgen

import (
	"fmt"
	"github.com/smartlet/wsdl2go/builtin"
	"go/format"
	"io"
)

func generateDefinitions(c *Context, pack string, out io.Writer, buf *Buffer) {

	buf.Line("package %s", pack)
	buf.Line("import (")
	buf.Line(`"context"`)
	buf.Line(`"encoding/xml"`)
	buf.Line(")\n")

	generateBuiltinType(c, buf)
	fmt.Fprintln(buf)
	generateInnerSimpleType(c, buf)
	fmt.Fprintln(buf)
	generateNamedSimpleType(c, buf)
	fmt.Fprintln(buf)
	generateInnerComplexType(c, buf)
	fmt.Fprintln(buf)
	generateNamedComplexType(c, buf)
	fmt.Fprintln(buf)
	generateNamedMessage(c, buf)
	fmt.Fprintln(buf)
	generatePortTypeInterface(c, buf)
	fmt.Fprintln(buf)
	generateBindingImplement(c, buf)
	fmt.Fprintln(buf)
	generateEnvelopeTypes(c, buf)
	data, err := format.Source(buf.Bytes())
	if err != nil {
		out.Write(buf.Bytes())
	} else {
		out.Write(data)
	}
}

func generateBuiltinType(c *Context, buf *Buffer) {
	buf.Line(builtin.Export(""))
	fmt.Fprintln(buf)
}

func generateSimpleType(c *Context, buf *Buffer, ts []*SimpleType) {
	for _, t := range ts {
		if t.deprecated {
			continue
		}
		gname := Identifier(t.Name)
		if t.Base != nil {
			buf.Line("type %s %v\n", gname, TypeName(t.Base))
			// TODO: Validate()
			if len(t.Enumeration) > 0 {
				buf.Line("const (")
				for _, v := range t.Enumeration {
					buf.Line("%v%v = %q", gname, Identifier(v), v)
				}
				buf.Line(")")
			}
		} else if t.List != nil {
			buf.Line("type %s []%v\n", gname, TypeName(t.List))

		} else if t.Union != nil {
			types := NewBuffer(128)
			for i, v := range t.Union {
				if i > 0 {
					types.WriteByte('|')
				}
				types.WriteString(TypeName(v))
			}
			buf.Line("type %s any // union(%s)", gname, types)
		} else {
			panic("invalid simpleType base")
		}

	}
}

func generateInnerSimpleType(c *Context, buf *Buffer) {
	for _, ts := range c.innerSimpleTypes.AllByNs() {
		generateSimpleType(c, buf, ts)
	}
}

func generateNamedSimpleType(c *Context, buf *Buffer) {
	for _, set := range c.namedSimpleTypes.All() {
		generateSimpleType(c, buf, set.All())
	}
}

func generateComplexType(c *Context, buf *Buffer, ts []*ComplexType) {

	for _, t := range ts {
		if t.deprecated {
			continue
		}

		buf.Line("type %s struct {", Identifier(t.Name))
		if t.Base != nil {
			switch bt := t.Base.(type) {
			case nil:
			case BuiltinType:
				buf.Line("CharData %s `xml:\",chardata,omitempty\"`", PointerTypeName(bt))
			case *SimpleType:
				buf.Line("CharData %s `xml:\",chardata,omitempty\"`", PointerTypeName(bt))
			case *ComplexType:
				buf.Line("%s `xml:\",omitempty\"`", TypeName(bt.Name)) // 不需要指针"*"形式!
			default:
				panic("invalid base type")
			}
		}
		for _, a := range t.Attributes {
			buf.Line("%s %s `xml:\"%s,attr,omitempty\"`", Identifier(a.Name), TypeName(a.Type), a.Name) // attribute默认不带前缀
		}
		for _, e := range t.Elements {
			var gtype string
			if e.MaxOccurs != "" && e.MaxOccurs != "0" && e.MaxOccurs != "1" {
				gtype += "[]"
			}
			gtype += PointerTypeName(e.Type)
			buf.Line("%s %s `xml:\"%s,omitempty\"`", Identifier(e.Name), gtype, c.QName(e.Ns, e.Name))
		}
		buf.Line("}\n")
	}
}

func generateInnerComplexType(c *Context, buf *Buffer) {
	for _, ts := range c.innerComplexTypes.AllByNs() {
		generateComplexType(c, buf, ts)
	}
}

func generateNamedComplexType(c *Context, buf *Buffer) {
	for _, set := range c.namedComplexTypes.All() {
		generateComplexType(c, buf, set.All())
	}
}

func generateNamedMessage(c *Context, buf *Buffer) {
	for _, ms := range c.namedMessages.All() {
		for _, m := range ms.All() {
			buf.Line("type %s struct {", Identifier(m.Name))
			for _, e := range m.Parts.All() {
				var gtype string
				if e.MaxOccurs != "" && e.MaxOccurs != "0" && e.MaxOccurs != "1" {
					gtype += "[]"
				}
				gtype += PointerTypeName(e.Type)
				buf.Line("%s %s `xml:\"%s,omitempty\"`", Identifier(e.Name), gtype, c.QName(e.Ns, e.Name))
			}
			buf.Line("}\n")
		}
	}
}

func generatePortTypeInterface(c *Context, buf *Buffer) {

	for _, pt := range c.namedPortTypes.All() {
		buf.Line("type %s interface {", Identifier(pt.Name))
		for _, op := range pt.Operations.All() {
			// 注意: Message不能使用TypeName()
			buf.Line("%s(ctx context.Context, in *%s) (*%s, error)", Identifier(op.Name), Identifier(op.Input.Name), Identifier(op.Output.Name))
		}
		buf.Line("}\n")
	}

}

func generateBindingImplement(c *Context, buf *Buffer) {
	for _, bd := range c.namedBindings.All() {
		buf.Line("type %s struct {", Identifier(bd.Name))
		buf.Line("client ")
		buf.Line("}\n")
		for _, op := range bd.Operations.All() {

		}
	}
}

func generateEnvelopeTypes(c *Context, buf *Buffer) {

}
