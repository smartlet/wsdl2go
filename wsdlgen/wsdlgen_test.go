package wsdlgen

import (
	"os"
	"testing"
)

const (
	wsdlFile   = "../test/services.wsdl"
	outputFile = `E:\temp\services.wsdl.go`
)

var DefaultPrefix = map[string]string{
	"http://schemas.xmlsoap.org/soap/envelope/":                    "s",
	"http://schemas.microsoft.com/exchange/services/2006/messages": "m",
	"http://schemas.microsoft.com/exchange/services/2006/types":    "t",
}

func TestGenerate(t *testing.T) {

	out, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	WsdlGen(wsdlFile, DefaultPrefix, "wsdl", out)

}
