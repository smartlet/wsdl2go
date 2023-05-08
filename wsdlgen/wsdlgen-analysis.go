package wsdlgen

import (
	"github.com/smartletn/wsdl2go/wsdl"
	"strconv"
	"strings"
)

func analysisDefinitions(c *Context) {

	c.trace("analysisDefinitions...")

	var ns, name string
	for _, op := range c.definitions.PortType.Operations {
		ns, name = c.definitions.QName(op.Input.Message)
		analysisNamedMessage(c, ns, name)
		ns, name = c.definitions.QName(op.Output.Message)
		analysisNamedMessage(c, ns, name)
	}
}

func analysisNamedMessage(c *Context, ns, name string) *Message {
	c.trace("analysisMessage: %v", name)

	for _, m := range c.definitions.Messages {
		if m.Name == name {
			rt := c.namedMessages.Get(ns, name)
			if rt != nil {
				return rt
			}
			rt = &Message{Ns: ns, Name: name}
			if c.namedMessages.Set(ns, name, rt) {
				panic("duplicate message:" + name)
			}
			rt.Parts = NewNamedSlice[*Element]()
			for _, p := range m.Parts {
				ns, name = c.definitions.QName(p.Element)
				rt.Parts.Set(p.Name, analysisNamedElement(c, ns, name))
			}

			return rt
		}
	}
	panic("missing message: " + name)
}

func analysisNamedElement(c *Context, ns, name string) *Element {
	c.trace("analysisElement: %v", name)

	sc := c.schemas.Get(ns)
	for _, e := range sc.Elements {
		if e.Name == name {
			rt := c.namedElements.Get(ns, name)
			if rt != nil {
				return rt
			}
			rt = &Element{Ns: ns, Name: name}
			if c.namedElements.Set(ns, name, rt) {
				panic("duplicate element: " + name)
			}
			return processElement(c, sc, e, rt, "")
		}
	}

	panic("missing element: " + name)
}

func analysisNamedType(c *Context, ns, name string) Type {

	c.trace("analysisNamedType: %v", name)

	if rt := Builtin(ns, name); rt != "" {
		return rt
	}
	if rt := c.namedSimpleTypes.Get(ns, name); rt != nil {
		return rt
	}
	if rt := c.namedComplexTypes.Get(ns, name); rt != nil {
		return rt
	}

	sc := c.schemas.Get(ns)

	for _, i := range sc.SimpleTypes {
		if name == i.Name {
			rt := &SimpleType{Ns: ns, Name: name}
			if c.namedSimpleTypes.Set(ns, name, rt) {
				panic("duplicate simpleType: " + name)
			}
			return processSimpleType(c, sc, i, rt)
		}
	}

	for _, i := range sc.ComplexTypes {
		if name == i.Name {
			rt := &ComplexType{Ns: ns, Name: name}
			if c.namedComplexTypes.Set(ns, name, rt) {
				panic("duplicate complexType: " + name)
			}
			return processComplexType(c, sc, i, rt)
		}
	}

	panic("invalid type: %v" + name)
}

func processSimpleType(c *Context, sc *wsdl.Schema, st *wsdl.SimpleType, rt *SimpleType) *SimpleType {

	c.trace("processSimpleType: %v", rt.Name)

	if rt.Name == "" {
		panic("require name")
	}

	if st.Restriction != nil {
		if st.Restriction.Base != "" {
			ns, name := sc.QName(st.Restriction.Base)
			rt.Base = analysisNamedType(c, ns, name)
		} else if st.Restriction.SimpleType != nil {
			it := &SimpleType{Ns: sc.TargetNamespace, Name: rt.Name + "Base"}
			c.innerSimpleTypes.Add(it.Ns, it.Name, it)
			rt.Base = processSimpleType(c, sc, st.Restriction.SimpleType, it)
		} else {
			panic("invalid restriction")
		}
		RestrictionSimpleType(st.Restriction, rt)
	} else if st.List != nil {
		if st.List.ItemType != "" {
			ns, name := sc.QName(st.List.ItemType)
			rt.List = analysisNamedType(c, ns, name)
		} else if st.List.SimpleType != nil {
			it := &SimpleType{Ns: sc.TargetNamespace, Name: rt.Name + "Item"}
			c.innerSimpleTypes.Add(it.Ns, it.Name, it)
			rt.List = processSimpleType(c, sc, st.List.SimpleType, it)
		} else {
			panic("invalid restriction")
		}
	} else if st.Union != nil {
		if st.Union.MemberTypes != "" {
			for _, i := range whitespace.Split(st.Union.MemberTypes, -1) {
				ns, name := sc.QName(i)
				rt.Union = append(rt.Union, analysisNamedType(c, ns, name))
			}
		}
		if len(st.Union.SimpleTypes) > 0 {
			for i, t := range st.Union.SimpleTypes {
				it := &SimpleType{Ns: sc.TargetNamespace, Name: rt.Name + "Union" + strconv.Itoa(i)}
				c.innerSimpleTypes.Add(it.Ns, it.Name, it)
				rt.Union = append(rt.Union, processSimpleType(c, sc, t, it))
			}
		}
	} else {
		panic("invalid type: " + rt.Name)
	}
	return rt
}

func processComplexType(c *Context, sc *wsdl.Schema, ct *wsdl.ComplexType, rt *ComplexType) *ComplexType {
	c.trace("processComplexType: %v", rt.Name)

	if rt.Name == "" {
		panic("require name")
	}

	if ct.SimpleContent != nil {
		processContent(c, sc, ct, ct.SimpleContent.Restriction, ct.SimpleContent.Extension, rt)
	} else if ct.ComplexContent != nil {
		processContent(c, sc, ct, ct.ComplexContent.Restriction, ct.ComplexContent.Extension, rt)
	} else {
		if ct.Group != nil {
			processGroup(c, sc, ct, ct.Group, rt, "")
		}
		if ct.Choice != nil {
			processChoice(c, sc, ct, ct.Choice, rt, "")
		}
		if ct.Sequence != nil {
			processSequence(c, sc, ct, ct.Sequence, rt, "")
		}
		if len(ct.Attributes) > 0 {
			for _, i := range ct.Attributes {
				processAttribute(c, sc, ct, i, rt)
			}
		}
		if len(ct.AttributeGroups) > 0 {
			for _, i := range ct.AttributeGroups {
				processAttributeGroup(c, sc, ct, i, rt)
			}
		}
	} // 可能是abstract

	return rt
}

func processContent(c *Context, sc *wsdl.Schema, ct *wsdl.ComplexType, re *wsdl.Restriction, ex *wsdl.Extension, rt *ComplexType) {
	c.trace("processContent...")

	if re != nil {
		if re.Base != "" {
			ns, name := sc.QName(re.Base)
			rt.Base = analysisNamedType(c, ns, name)
		} else if re.SimpleType != nil {
			st := &SimpleType{Ns: sc.TargetNamespace, Name: ct.Name + "Type"}
			c.innerSimpleTypes.Add(st.Ns, st.Name, st)
			rt.Base = processSimpleType(c, sc, re.SimpleType, st)
		} else {
			panic("invalid type...")
		}
		RestrictionComplexType(re, rt)
	} else if ex != nil {
		if ex.Base != "" {
			ns, name := sc.QName(ex.Base)
			rt.Base = analysisNamedType(c, ns, name)
		}
		if ex.Group != nil {
			processGroup(c, sc, ct, ex.Group, rt, "")
		}
		if ex.Choice != nil {
			processChoice(c, sc, ct, ex.Choice, rt, "")
		}
		if ex.Sequence != nil {
			processSequence(c, sc, ct, ex.Sequence, rt, "")
		}
		if len(ex.Attributes) > 0 {
			for _, i := range ex.Attributes {
				processAttribute(c, sc, ct, i, rt)
			}
		}
		if len(ex.AttributeGroups) > 0 {
			for _, i := range ex.AttributeGroups {
				processAttributeGroup(c, sc, ct, i, rt)
			}
		}
	} else {
		panic("invalid content: " + ct.Name)
	}

}

func processChoice(c *Context, sc *wsdl.Schema, ct *wsdl.ComplexType, gp *wsdl.Choice, rt *ComplexType, maxOccurs string) {
	c.trace("processChoice...")

	maxOccurs = nvl(gp.MaxOccurs, maxOccurs)
	if len(gp.Elements) > 0 {
		for _, i := range gp.Elements {
			rt.Elements = append(rt.Elements, processElement(c, sc, i, &Element{Ns: sc.TargetNamespace, Name: i.Name}, maxOccurs))
		}
	}
	if len(gp.Groups) > 0 {
		for _, i := range gp.Groups {
			processGroup(c, sc, ct, i, rt, maxOccurs)
		}
	}
	if len(gp.Choices) > 0 {
		for _, i := range gp.Choices {
			processChoice(c, sc, ct, i, rt, maxOccurs)
		}
	}
	if len(gp.Sequences) > 0 {
		for _, i := range gp.Sequences {
			processSequence(c, sc, ct, i, rt, maxOccurs)
		}
	}
}

func processSequence(c *Context, sc *wsdl.Schema, ct *wsdl.ComplexType, gp *wsdl.Sequence, rt *ComplexType, maxOccurs string) {
	c.trace("processSequence...")
	// 处理重复
	maxOccurs = nvl(gp.MaxOccurs, maxOccurs)
	if len(gp.Elements) > 0 {
		for _, i := range gp.Elements {
			rt.Elements = append(rt.Elements, processElement(c, sc, i, &Element{Ns: sc.TargetNamespace, Name: i.Name}, maxOccurs))
		}
	}
	if len(gp.Groups) > 0 {
		for _, i := range gp.Groups {
			processGroup(c, sc, ct, i, rt, maxOccurs)
		}
	}
	if len(gp.Choices) > 0 {
		for _, i := range gp.Choices {
			processChoice(c, sc, ct, i, rt, maxOccurs)
		}
	}
	if len(gp.Sequences) > 0 {
		for _, i := range gp.Sequences {
			processSequence(c, sc, ct, i, rt, maxOccurs)
		}
	}
}

func processElement(c *Context, sc *wsdl.Schema, e *wsdl.Element, rt *Element, maxOccurs string) *Element {

	c.trace("processElement: %v", nvl(e.Name, e.Ref))

	if e.Ref != "" {
		ns, name := sc.QName(e.Ref)
		sc = c.schemas.Get(ns)
		for _, i := range sc.Elements {
			if name == i.Name {
				e = i
				break
			}
		}
		if e.Name != name {
			panic(" invalid ref: " + e.Ref)
		}
		// FIXBUG: ref
		rt.Ns = ns
		rt.Name = e.Name
	}

	rt.Default = e.Default
	rt.Fixed = e.Fixed
	rt.Use = e.Use
	rt.MaxOccurs = nvl(e.MaxOccurs, maxOccurs)
	if e.Type != "" {
		ns, name := sc.QName(e.Type)
		rt.Type = analysisNamedType(c, ns, name)
	} else if e.SimpleType != nil {
		it := &SimpleType{Ns: sc.TargetNamespace, Name: e.Name + "Type"}
		c.innerSimpleTypes.Add(it.Ns, it.Name, it)
		rt.Type = processSimpleType(c, sc, e.SimpleType, it)
	} else if e.ComplexType != nil {
		it := &ComplexType{Ns: sc.TargetNamespace, Name: e.Name + "Type"}
		c.innerComplexTypes.Add(it.Ns, it.Name, it)
		rt.Type = processComplexType(c, sc, e.ComplexType, it)
	} else {
		panic("invalid type: " + rt.Name)
	}
	return rt
}

func processGroup(c *Context, sc *wsdl.Schema, ct *wsdl.ComplexType, gp *wsdl.Group, rt *ComplexType, maxOccurs string) {

	c.trace("processElement: %v", nvl(gp.Name, gp.Ref))

	if gp.Ref != "" {
		ns, name := sc.QName(gp.Ref)
		sc = c.schemas.Get(ns)
		for _, i := range sc.Groups {
			if name == i.Name {
				gp = i
				break
			}
		}
		if gp.Name != name {
			panic("invalid ref: " + gp.Ref)
		}
	}
	// 处理重复
	maxOccurs = nvl(gp.MaxOccurs, maxOccurs)
	if gp.Choice != nil {
		processChoice(c, sc, ct, gp.Choice, rt, maxOccurs)
	} else if gp.Sequence != nil {
		processSequence(c, sc, ct, gp.Sequence, rt, maxOccurs)
	} else {
		panic("invalid group: " + gp.Name)
	}

}

func processAttribute(c *Context, sc *wsdl.Schema, ct *wsdl.ComplexType, a *wsdl.Attribute, rt *ComplexType) {
	c.trace("processAttribute: %v", nvl(a.Name, a.Ref))

	// 特殊属性
	if a.Ref == "xml:lang" || a.Ref == "xsi:nil" {
		idx := strings.IndexByte(a.Ref, ':')
		rt.Attributes = append(rt.Attributes, &Attribute{
			Ns:   a.Ref[:idx],
			Name: a.Ref[idx+1:],
			Type: BuiltinType("XsString"),
		})
		return
	}

	if a.Ref != "" {
		ns, name := sc.QName(a.Ref)
		sc = c.schemas.Get(ns)
		for _, i := range sc.Attributes {
			if name == i.Name {
				a = i
				break
			}
		}
		if a.Name != name {
			panic("invalid ref: " + a.Ref)
		}
	}

	at := &Attribute{Ns: sc.TargetNamespace, Name: a.Name}
	at.Default = a.Default
	at.Fixed = a.Fixed
	at.Use = a.Use
	if a.Type != "" {
		ns, name := sc.QName(a.Type)
		at.Type = analysisNamedType(c, ns, name)
	} else if a.SimpleType != nil {
		st := &SimpleType{Ns: sc.TargetNamespace, Name: a.Name + "Type"}
		c.innerSimpleTypes.Add(st.Ns, st.Name, st)
		at.Type = processSimpleType(c, sc, a.SimpleType, st)
	} else {
		panic("invalid type: " + a.Name)
	}
	rt.Attributes = append(rt.Attributes, at)
}

func processAttributeGroup(c *Context, sc *wsdl.Schema, ct *wsdl.ComplexType, ag *wsdl.AttributeGroup, rt *ComplexType) {
	c.trace("processAttributeGroup: %v", nvl(ag.Name, ag.Ref))
	if ag.Ref != "" {
		ns, name := sc.QName(ag.Ref)
		sc = c.schemas.Get(ns)
		for _, i := range sc.AttributeGroups {
			if name == i.Name {
				ag = i
				break
			}
		}
		if ag.Name != name {
			panic("invalid ref: " + ag.Ref)
		}
	}

	for _, i := range ag.Attributes {
		processAttribute(c, sc, ct, i, rt)
	}
	for _, i := range ag.AttributeGroups {
		processAttributeGroup(c, sc, ct, i, rt)
	}
}
