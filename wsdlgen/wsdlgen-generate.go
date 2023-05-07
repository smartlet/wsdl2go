package wsdlgen

import (
	"fmt"
	"github.com/smartletn/wsdlgen/builtin"
)

func generateBuiltinType(c *Context, buf *Buffer) {
	buf.Write(builtin.Export(""))
	buf.WriteByte('\n')
}

func generateSimpleType(c *Context, buf *Buffer, ts []*SimpleType) {
	for _, t := range ts {
		gname := Identifier(t.Name)
		if t.Base != nil {
			fmt.Fprintf(buf, "type %s %v\n\n", gname, TypeName(t.Base))
			// TODO: Validate()
			if len(t.Enumeration) > 0 {
				fmt.Fprintf(buf, "const (\n")
				for _, v := range t.Enumeration {
					fmt.Fprintf(buf, "%v%v = %q\n", gname, Identifier(v), v)
				}
				fmt.Fprintf(buf, ")\n")
			}
		} else if t.List != nil {
			fmt.Fprintf(buf, "type %s []%v\n\n", gname, TypeName(t.List))

		} else if t.Union != nil {
			types := NewBuffer(128)
			for i, v := range t.Union {
				if i > 0 {
					types.WriteByte('|')
				}
				types.WriteString(TypeName(v))
			}
			fmt.Fprintf(buf, "type %s any // union(%s)\n\n", gname, types)
		} else {
			panic("invalid simpleType base")
		}

	}
}

func generateInnerSimpleType(c *Context, buf *Buffer) {
	for _, ts := range c.innerSimpleTypes.All() {
		generateSimpleType(c, buf, ts)
	}
}

func generateNamedSimpleType(c *Context, buf *Buffer) {
	for _, set := range c.namedSimpleTypes.All() {
		generateSimpleType(c, buf, set.All())
	}
}

func generateInnerComplexType(c *Context, buf *Buffer) {

}

func generateNamedComplexType(c *Context, buf *Buffer) {

}

func generateNamedElement(c *Context, buf *Buffer) {

}

func generateNamedMessage(c *Context, buf *Buffer) {

}

func generatePortTypeInterface(c *Context, buf *Buffer) {

}

func generateBindingImplement(c *Context, buf *Buffer) {

}
