package wsdl

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"io"
)

func DecodeDefinitions(r io.Reader) (*Definitions, error) {
	var d Definitions
	decoder := xml.NewDecoder(r)
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(&d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func DecodeDeSchema(r io.Reader) (*Schema, error) {
	var d Schema
	decoder := xml.NewDecoder(r)
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(&d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
