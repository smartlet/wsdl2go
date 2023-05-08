package wsdlgen

import (
	"os"
	"testing"
)

const (
	wsdl     = "../test/services.wsdl"
	request  = `E:\temp\request.go`
	response = `E:\temp\response.go`
)

func TestGenerate(t *testing.T) {

	req, err := os.Create(request)
	if err != nil {
		panic(err)
	}
	defer req.Close()

	rsp, err := os.Create(response)
	if err != nil {
		panic(err)
	}
	defer rsp.Close()

	WsdlGen(wsdl, Output{
		ResponsePackage: "response",
		ResponseSources: rsp,
		RequestPackage:  "request",
		RequestSources:  rsp,
	})

}
