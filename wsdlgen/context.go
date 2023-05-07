package wsdlgen

import (
	"fmt"
	"github.com/smartletn/wsdl2go/wsdl"
	"io"
	"os"
	"path/filepath"
)

func NewContext(path string) *Context {
	ctx := &Context{
		wsdl:              path,
		base:              filepath.Dir(path),
		schemas:           NewNamedSlice[*wsdl.Schema](),
		namedSimpleTypes:  NewNsNamedSlice[*SimpleType](),
		innerSimpleTypes:  NewNsOrderSlice[*SimpleType](),
		namedComplexTypes: NewNsNamedSlice[*ComplexType](),
		innerComplexTypes: NewNsOrderSlice[*ComplexType](),
		namedElements:     NewNsNamedSlice[*Element](),
		namedMessages:     NewNsNamedSlice[*Message](),
		prefixes:          make(map[string]string),
	}

	// 添加输出时的默认的前缀
	for k, v := range defaultPrefix {
		ctx.prefixes[k] = v
	}

	return ctx
}

type Context struct {
	wsdl              string
	base              string                      // 目录. 用于相对路径的处理.
	schemas           *NamedSlice[*wsdl.Schema]   // 根据targetNamespace索引
	namedSimpleTypes  *NsNamedSlice[*SimpleType]  // 顶层带有name可以被ref的simpleType
	innerSimpleTypes  *NsOrderSlice[*SimpleType]  // 在element/attribute/base内声明的无name的simpleType
	namedComplexTypes *NsNamedSlice[*ComplexType] // 顶层带有name可以被ref的complexType
	innerComplexTypes *NsOrderSlice[*ComplexType] // 在element/attribute/base内声明的无name的complexType
	namedElements     *NsNamedSlice[*Element]     // 顶层带有name可以被ref的element
	namedMessages     *NsNamedSlice[*Message]     // 顶层带有name可以被ref的message
	definitions       *wsdl.Definitions
	prefixes          map[string]string
	traceWriter       io.WriteCloser
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
