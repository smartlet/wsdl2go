package wsdlgen

import "strings"

// 特殊处理场景. 匿名相同则加上ns后缀方便区别
func specialsDefinitions(c *Context) {
	for _, ks := range c.innerSimpleTypes.AllByKey() {
		if ln := len(ks); ln > 1 {
			// 去重
			sn := ln
			for i := 1; i < ln; i++ {
				for j := 0; j < i; j++ {
					if !ks[j].deprecated && equalsType(ks[j], ks[i]) {
						ks[i].deprecated = true
						sn--
						break
					}
				}
			}
			if sn > 1 {
				for _, k := range ks {
					k.Name += strings.ToUpper(c.prefixes[k.Ns])
				}
			}
		}
	}
	for _, ks := range c.innerComplexTypes.AllByKey() {
		if ln := len(ks); ln > 1 {
			// 去重
			sn := ln
			for i := 1; i < ln; i++ {
				for j := 0; j < i; j++ {
					if !ks[j].deprecated && equalsType(ks[j], ks[i]) {
						ks[i].deprecated = true
						sn--
						break
					}
				}
			}
			if sn > 1 {
				for _, k := range ks {
					k.Name += strings.ToUpper(c.prefixes[k.Ns])
				}
			}
		}
	}
}

func equalsType(t Type, v Type) bool {

	if t == v {
		return true
	}

	switch t := t.(type) {
	case nil:
		return v == nil
	case BuiltinType:
		bv, ok := v.(BuiltinType)
		if !ok {
			return false
		}
		return t == bv

	case *SimpleType:
		sv, ok := v.(*SimpleType)
		if !ok {
			return false
		}

		if t.Ns != sv.Ns {
			return false
		}
		if t.Name != sv.Name {
			return false
		}

		if !equalsType(t.Base, sv.Base) {
			return false
		}

		if !equalsStrings(t.MinExclusive, sv.MinExclusive) {
			return false
		}

		if !equalsStrings(t.MinInclusive, sv.MinInclusive) {
			return false
		}

		if !equalsStrings(t.MaxExclusive, sv.MaxExclusive) {
			return false
		}

		if !equalsStrings(t.MaxInclusive, sv.MaxInclusive) {
			return false
		}

		if !equalsStrings(t.TotalDigits, sv.TotalDigits) {
			return false
		}

		if !equalsStrings(t.FractionDigits, sv.FractionDigits) {
			return false
		}

		if !equalsStrings(t.Length, sv.Length) {
			return false
		}

		if !equalsStrings(t.MinLength, sv.MinLength) {
			return false
		}

		if !equalsStrings(t.MaxLength, sv.MaxLength) {
			return false
		}

		if !equalsStrings(t.WhiteSpace, sv.WhiteSpace) {
			return false
		}

		if !equalsStrings(t.Pattern, sv.Pattern) {
			return false
		}

		if !equalsStrings(t.Pattern, sv.Pattern) {
			return false
		}

		if !equalsStrings(t.Enumeration, sv.Enumeration) {
			return false
		}

		if !equalsType(t.List, sv.List) {
			return false
		}

		if tn, sn := len(t.Union), len(sv.Union); tn != sn {
			return false
		} else {
			for i := 0; i < tn; i++ {
				if !equalsType(t.Union[i], sv.Union[i]) {
					return false
				}
			}
		}

		return true

	case *ComplexType:
		sv, ok := v.(*ComplexType)
		if !ok {
			return false
		}

		if t.Ns != sv.Ns {
			return false
		}
		if t.Name != sv.Name {
			return false
		}

		if !equalsType(t.Base, sv.Base) {
			return false
		}

		if !equalsStrings(t.MinExclusive, sv.MinExclusive) {
			return false
		}

		if !equalsStrings(t.MinInclusive, sv.MinInclusive) {
			return false
		}

		if !equalsStrings(t.MaxExclusive, sv.MaxExclusive) {
			return false
		}

		if !equalsStrings(t.MaxInclusive, sv.MaxInclusive) {
			return false
		}

		if !equalsStrings(t.TotalDigits, sv.TotalDigits) {
			return false
		}

		if !equalsStrings(t.FractionDigits, sv.FractionDigits) {
			return false
		}

		if !equalsStrings(t.Length, sv.Length) {
			return false
		}

		if !equalsStrings(t.MinLength, sv.MinLength) {
			return false
		}

		if !equalsStrings(t.MaxLength, sv.MaxLength) {
			return false
		}

		if !equalsStrings(t.WhiteSpace, sv.WhiteSpace) {
			return false
		}

		if !equalsStrings(t.Pattern, sv.Pattern) {
			return false
		}

		if !equalsStrings(t.Pattern, sv.Pattern) {
			return false
		}

		if !equalsStrings(t.Enumeration, sv.Enumeration) {
			return false
		}

		if !equalsAttributes(t.Attributes, sv.Attributes) {
			return false
		}

		if !equalsElements(t.Elements, sv.Elements) {
			return false
		}
		return true
	}
	panic("invalid type")
}

func equalsStrings(v1 []string, v2 []string) bool {
	n1, n2 := len(v1), len(v2)
	if n1 != n2 {
		return false
	}
	for i := 0; i < n1; i++ {
		if v1[i] != v2[i] {
			return false
		}
	}
	return true
}

func equalsAttributes(as1 []*Attribute, as2 []*Attribute) bool {
	an1, an2 := len(as1), len(as2)
	if an1 != an2 {
		return false
	}
	for i := 0; i < an1; i++ {
		if !equalsAttribute(as1[i], as2[i]) {
			return false
		}
	}
	return true
}

func equalsAttribute(a1 *Attribute, a2 *Attribute) bool {
	if a1 == a2 {
		return true
	}
	if a1.Ns != a2.Ns {
		return false
	}
	if a1.Name != a2.Name {
		return false
	}
	if a1.Default != a2.Default {
		return false
	}
	if a1.Fixed != a2.Fixed {
		return false
	}
	if a1.Use != a2.Use {
		return false
	}
	if !equalsType(a1.Type, a2.Type) {
		return false
	}
	return true
}

func equalsElements(as1 []*Element, as2 []*Element) bool {
	an1, an2 := len(as1), len(as2)
	if an1 != an2 {
		return false
	}
	for i := 0; i < an1; i++ {
		if !equalsElement(as1[i], as2[i]) {
			return false
		}
	}
	return true
}

func equalsElement(a1 *Element, a2 *Element) bool {
	if a1 == a2 {
		return true
	}
	if a1.Ns != a2.Ns {
		return false
	}
	if a1.Name != a2.Name {
		return false
	}
	if a1.Default != a2.Default {
		return false
	}
	if a1.Fixed != a2.Fixed {
		return false
	}
	if a1.Use != a2.Use {
		return false
	}
	if a1.MaxOccurs != a2.MaxOccurs {
		return false
	}
	if !equalsType(a1.Type, a2.Type) {
		return false
	}
	return true
}
