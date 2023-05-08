package wsdlgen

import (
	"io"
	"path/filepath"
)

type Output struct {
	RequestPackage  string
	RequestSources  io.Writer
	ResponsePackage string
	ResponseSources io.Writer
}

func WsdlGen(wsdl string, output Output) {

	abs, err := filepath.Abs(wsdl)
	if err != nil {
		panic(err)
	}

	ctx := NewContext(abs)
	defer ctx.close()

	decodeDefinitions(ctx)
	analysisDefinitions(ctx)
	compressDefinitions(ctx)

	buf := NewBuffer(2048)

	buf.Reset()
	generateDefinitions(ctx, output.ResponsePackage, output.ResponseSources, buf)

	ctx.Prefixes(defaultPrefix)
	buf.Reset()
	generateDefinitions(ctx, output.RequestPackage, output.RequestSources, buf)

}
