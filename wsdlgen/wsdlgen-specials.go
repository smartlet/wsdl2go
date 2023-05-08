package wsdlgen

import "strings"

// 特殊处理场景. 匿名相同则加上ns后缀方便区别
func specialsDefinitions(c *Context) {
	for _, ks := range c.innerSimpleTypes.AllByKey() {
		if len(ks) > 1 {
			for _, k := range ks {
				k.Name = strings.ToUpper(c.prefixes[k.Ns]) + k.Name
			}
		}
	}
	for _, ks := range c.innerComplexTypes.AllByKey() {
		if len(ks) > 1 {
			for _, k := range ks {
				k.Name = strings.ToUpper(c.prefixes[k.Ns]) + k.Name
			}
		}
	}
}
