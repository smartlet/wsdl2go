package wsdlgen

import (
	"os"
	"testing"
)

const (
	servicesWsdlFile = "../test/services.wsdl"
	servicesGoFile   = `E:\temp\generate.go`
)

func TestGenerate(t *testing.T) {

	out, err := os.Create(servicesGoFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	WsdlGen(servicesWsdlFile, "ews", out)

}
