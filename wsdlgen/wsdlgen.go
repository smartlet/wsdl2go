package wsdlgen

import (
	"io"
	"path/filepath"
)

func WsdlGen(wsdl string, prefixes map[string]string, pack string, out io.Writer) {

	abs, err := filepath.Abs(wsdl)
	if err != nil {
		panic(err)
	}

	ctx := NewContext(abs, prefixes)
	defer ctx.close()

	decodeDefinitions(ctx)
	analysisDefinitions(ctx)
	compressDefinitions(ctx)

	buf := NewBuffer(2048)

	ctx.Prefixes(prefixes)
	generateDefinitions(ctx, pack, out, buf)

}
