package wsdl

import (
	"fmt"
	"os"
	"testing"
)

func TestDecodeDefinitions(t *testing.T) {
	file, err := os.Open("../wsdl/services.wsdl")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ds, err := DecodeDefinitions(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(ds.Namespaces)
}

func TestDecodeDeSchema(t *testing.T) {
	file, err := os.Open("../wsdl/types.xsd")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ts, err := DecodeDeSchema(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(ts.Namespaces)
}
