package wsdlgen

import (
	"github.com/smartlet/wsdl2go/builtin"
	"os"
	"testing"
)

const (
	wsdlFile   = "../test/services.wsdl"
	outputFile = `E:\temp\services.wsdl.go`
)

func TestGenerate(t *testing.T) {

	out, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	WsdlGen(wsdlFile, builtin.XmlnsPrefix, "wsdl", out)

}
