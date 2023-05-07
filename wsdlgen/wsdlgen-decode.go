package wsdlgen

import (
	"github.com/smartletn/wsdlgen/wsdl"
	"os"
	"path/filepath"
)

func decodeDefinitions(c *Context) {

	c.trace("decodeDefinitions: %v", c.wsdl)

	file, err := os.Open(c.wsdl)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	c.definitions, err = wsdl.DecodeDefinitions(file)
	if err != nil {
		panic(err)
	}
	for _, imp := range c.definitions.Schema.Imports {
		if imp.SchemaLocation != "" {
			decodeSchema(c, imp.Namespace, imp.SchemaLocation)
		}
	}
}

func decodeSchema(c *Context, ns, schemaLocation string) {

	c.trace("decodeSchema: %v, %v", ns, schemaLocation)

	if !filepath.IsAbs(schemaLocation) {
		schemaLocation = filepath.Join(c.base, schemaLocation)
	}
	file, err := os.Open(schemaLocation)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ts, err := wsdl.DecodeDeSchema(file)
	if err != nil {
		panic(err)
	}
	if ns != ts.TargetNamespace {
		panic("targetNamespace invalid")
	}
	c.schemas.Set(ns, ts)
	for _, imp := range ts.Imports {
		if imp.SchemaLocation != "" {
			decodeSchema(c, imp.Namespace, imp.SchemaLocation)
		}
	}
}
