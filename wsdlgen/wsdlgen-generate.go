package wsdlgen

import (
	"fmt"
	"github.com/smartlet/wsdl2go/builtin"
	"go/format"
	"io"
)

func generateDefinitions(c *Context, pack string, out io.Writer, buf *Buffer) {
	generateBuiltinType(c, pack, buf)
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

func generateEnvelopeTypes(c *Context, buf *Buffer) {

}

func generateBuiltinType(c *Context, pack string, buf *Buffer) {
	buf.Line(builtin.Export(pack))
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
	/*
		- ComplexType: innerComplexType/namedComplexType

		    // 根据chardata/attribute/element可以分成3类:
		    // 1. 包含chardata/attribute, 对应simpleContent
		    // 2. 包含chardata/attribute/element, 对应complexContent
		    // 3. 包含attribute/element, 没有chardata. 对应除simpleContent/complexContent外的情况.
		    // complexType只支持上述3种形式.

		    // 在此基础上扩展只有chardata的情况. 则complexType的定义形式:
		    // 1. base + facets
		    // 2. attributes
		    // 3. elements
		    // 一个前提: complexType必须是struct结构. 否则会与前面描述的Element规则冲突.

		    // 若base.Type是BuiltinType或*SimpleType.则
		    type <ComplexType.Name> struct {
		        CharData <base.Type> `xml:",chardata"`
		    }
		    func (t *{ComplexType.Name})Validate() bool{
		        // 针对t.CharData进行校验....
		    }

		    // 若base.Type是*ComplexType. 根据前提: 其必是另一个struct
		    type <ComplexType.Name> struct {
		        <base.Type> // 采用nested struct的形式继承. 此处不需指针形式.
		    }

		    // 对于attributes. 输出为 {attribute.Name} {attribute.Type} `xml:"name,attr,omitempty"`形式
		    type <ComplexType.Name> struct {
		        // base处理
		        // attribute处理. 其中attribute.Type肯定是BuiltinType或SimpleType. 不需要指针"*"形式!
		        {attribute.Name} {attribute.Type} `xml:"name,attr,omitempty"`

		    }

		    // 对于elements. 处理方式见上...
	*/

	for _, t := range ts {
		if t.deprecated {
			continue
		}

		buf.Line("type %s struct {", Identifier(t.Name))
		if t.Base != nil {
			switch bt := t.Base.(type) {
			case nil:
			case BuiltinType:
				buf.Line("CharData %s `xml:\",chardata,omitempty\"`", FieldTypeName(bt))
			case *SimpleType:
				buf.Line("CharData %s `xml:\",chardata,omitempty\"`", FieldTypeName(bt))
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
			gtype += FieldTypeName(e.Type)
			buf.Line("%s %s `xml:\"%s,omitempty\"`", Identifier(e.Name), gtype, c.QName(e.Ns, e.Name))
		}
		buf.Line("}")
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
				gtype += FieldTypeName(e.Type)
				buf.Line("%s %s `xml:\"%s,omitempty\"`", Identifier(e.Name), gtype, c.QName(e.Ns, e.Name))
			}
			buf.Line("}")
		}
	}
}

func generatePortTypeInterface(c *Context, buf *Buffer) {

}

func generateBindingImplement(c *Context, buf *Buffer) {

}
