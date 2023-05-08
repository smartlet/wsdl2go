package wsdlgen

import (
	"fmt"
	"go/format"
	"io"
	"path/filepath"
)

func WsdlGen(wsdl string, pack string, out io.Writer) {

	abs, err := filepath.Abs(wsdl)
	if err != nil {
		panic(err)
	}

	ctx := NewContext(abs)
	defer ctx.close()

	decodeDefinitions(ctx)
	analysisDefinitions(ctx)
	specialsDefinitions(ctx)

	buf := NewBuffer(2048)
	// defer buffer clean

	fmt.Fprintln(buf, "package", pack)
	fmt.Fprintln(buf)
	generateBuiltinType(ctx, buf)
	fmt.Fprintln(buf)
	generateInnerSimpleType(ctx, buf)
	fmt.Fprintln(buf)
	generateNamedSimpleType(ctx, buf)
	fmt.Fprintln(buf)
	generateInnerComplexType(ctx, buf)
	fmt.Fprintln(buf)
	generateNamedComplexType(ctx, buf)
	fmt.Fprintln(buf)
	generateNamedElement(ctx, buf)
	fmt.Fprintln(buf)
	generateNamedMessage(ctx, buf)
	fmt.Fprintln(buf)
	generatePortTypeInterface(ctx, buf)
	fmt.Fprintln(buf)
	generateBindingImplement(ctx, buf)
	fmt.Fprintln(buf)
	data, err := format.Source(buf.Bytes())
	if err != nil {
		out.Write(buf.Bytes())
	} else {
		out.Write(data)
	}
}
