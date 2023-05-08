package wsdlgen

import (
	"github.com/smartletn/wsdl2go/builtin"
)

func generateBuiltinType(c *Context, buf *Buffer) {
	buf.Write(builtin.Export(""))
	buf.WriteByte('\n')
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
		# 对于ComplexType
		- 如果len(attributes)==0 && len(elements) == 0
		    - 如果base非空
		        形式为: type <ComplexType> <base>
		        在Validate()添加facets针对ComplexType本身.
		        此时退化为*SimpleType
		    - 如果base为空
		        形式为: type <ComplexType> any
		- 如果len(attributes) > 0 || len(elements) > 0
		    形式为: type <ComplexType> struct {
		        <base> 根据BuiltinType/*SimpleType或*ComplexType分别处理
		        <attributes> 渲染为`xml:"xxx,attr,omitempty"`
		        <elements> 渲染为`xml:"xxx,omitempty"`
		    }
		    - base是BuiltinType或*SimpleType.
		        Data <base> `xml:,chardata`
		        在Validate()添加facets针对ComplexType.Data
		    - base是*ComplexType则使用nested type
		在根据ComplexType渲染Element时, 也应该根据上述规则.
	*/

	for _, t := range ts {
		if t.deprecated {
			continue
		}
		gname := Identifier(t.Name)

		if len(t.Attributes) == 0 && len(t.Elements) == 0 {
			// 退化为SimleType处理
			if t.Base != nil {
				buf.Line("type %s %v\n", gname, TypeName(t.Base))
			} else {
				buf.Line("type %s any\n", gname)
			}
			// TODO: Validate()
			if len(t.Enumeration) > 0 {
				buf.Line("const (")
				for _, v := range t.Enumeration {
					buf.Line("%v%v = %q", gname, Identifier(v), v)
				}
				buf.Line(")")
			}
		} else {
			buf.Line("type %s struct {", gname)

			if t.Base != nil {
				if _, ok := t.Base.(*ComplexType); ok {
					buf.Line("%v", TypeName(t.Base)) // nested struct
				} else {
					buf.Line("CharData %v `xml:\",chardata,omitempty\"`", TypeName(t.Base))
				}
			}

			for _, a := range t.Attributes {
				buf.Line("%v %v `xml:\"%s,attr,omitempty\"`", Identifier(a.Name), FieldTypeName(a.Type), c.QName(a.Ns, a.Name))
			}

			for _, e := range t.Elements {
				buf.Line("%v %v `xml:\"%s,omitempty\"`", Identifier(e.Name), ElementFieldTypeName(e), c.QName(e.Ns, e.Name))
			}

			buf.Line("}\n")

			// TODO: Validate()
			if len(t.Enumeration) > 0 {
				// 理论不会有此值,但还是加上以防万一.
				buf.Line("const (")
				for _, v := range t.Enumeration {
					buf.Line("%v%v = %q", gname, Identifier(v), v)
				}
				buf.Line(")")
			}
		}
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

func generateNamedElement(c *Context, buf *Buffer) {

}

func generateNamedMessage(c *Context, buf *Buffer) {

}

func generatePortTypeInterface(c *Context, buf *Buffer) {

}

func generateBindingImplement(c *Context, buf *Buffer) {

}
