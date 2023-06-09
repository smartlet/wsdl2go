package wsdlgen

import (
	"fmt"
	"github.com/smartlet/wsdl2go/wsdl"
	"io"
	"os"
	"path/filepath"
)

func NewContext(path string, prefixes map[string]string) *Context {
	return &Context{
		wsdl:              path,
		base:              filepath.Dir(path),
		schemas:           NewNamedSlice[*wsdl.Schema](),
		namedSimpleTypes:  NewNsNamedSlice[*SimpleType](),
		innerSimpleTypes:  NewNsOrderSlice[*SimpleType](),
		namedComplexTypes: NewNsNamedSlice[*ComplexType](),
		innerComplexTypes: NewNsOrderSlice[*ComplexType](),
		namedMessages:     NewNsNamedSlice[*Message](),
		namedPortTypes:    NewNamedSlice[*PortType](),
		namedBindings:     NewNamedSlice[*Binding](),
		prefixes:          prefixes,
	}
}

type Context struct {
	wsdl              string
	base              string                    // 目录. 用于相对路径的处理.
	schemas           *NamedSlice[*wsdl.Schema] // 根据targetNamespace索引
	definitions       *wsdl.Definitions
	namedSimpleTypes  *NsNamedSlice[*SimpleType]  // 顶层带有name可以被ref的simpleType
	innerSimpleTypes  *NsOrderSlice[*SimpleType]  // 在element/attribute/base内声明的无name的simpleType
	namedComplexTypes *NsNamedSlice[*ComplexType] // 顶层带有name可以被ref的complexType
	innerComplexTypes *NsOrderSlice[*ComplexType] // 在element/attribute/base内声明的无name的complexType
	namedMessages     *NsNamedSlice[*Message]     // 顶层带有name可以被ref的message
	namedPortTypes    *NamedSlice[*PortType]
	namedBindings     *NamedSlice[*Binding]
	prefixes          map[string]string
	traceWriter       io.WriteCloser
}

func (c *Context) Prefixes(m map[string]string) {
	c.prefixes = m
}

// QName 如果设置了prefixes则采用p:name形式, 否则使用"namespace-URL name"形式
func (c *Context) QName(ns, name string) string {

	if ns == "" {
		return name
	}

	p := c.prefixes[ns]
	if p != "" {
		return p + ":" + name
	}
	return ns + " " + name

}

func (c *Context) trace(format string, args ...any) {
	if c.traceWriter == nil {
		c.traceWriter, _ = os.Create(`E:\temp\trace.log`)
	}
	fmt.Fprintf(c.traceWriter, format, args...)
	fmt.Fprintln(c.traceWriter)
}

func (c *Context) close() {
	if c.traceWriter != nil {
		c.traceWriter.Close()
	}
}
