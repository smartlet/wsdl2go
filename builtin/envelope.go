package builtin

import "encoding/xml"

const (
	XmlnsS = "http://schemas.xmlsoap.org/soap/envelope/"
	XmlnsT = "https://schemas.microsoft.com/exchange/services/2006/types"
	XmlnsM = "https://schemas.microsoft.com/exchange/services/2006/messages"
	XmlnsX = ""
)

type RequestEnvelope struct {
	XMLName xml.Name       `xml:"s:Envelope"`
	XmlnsS  string         `xml:"xmlns:s,attr,omitempty"`
	XmlnsT  string         `xml:"xmlns:t,attr,omitempty"`
	XmlnsM  string         `xml:"xmlns:m,attr,omitempty"`
	Header  *RequestHeader `xml:",omitempty"`
	Body    RequestBody    `xml:",omitempty"`
}

type RequestHeader struct {
	XMLName xml.Name `xml:"s:Header"`
	Headers []any    `xml:",omitempty"`
}

type RequestBody struct {
	XMLName xml.Name    `xml:"s:Body"`
	Content interface{} `xml:",omitempty"`
}

type ResponseEnvelope struct {
	XMLName xml.Name        `xml:"Envelope"`
	Header  *ResponseHeader `xml:",omitempty"`
	Body    ResponseBody    `xml:",omitempty"`
}

type ResponseHeader struct {
	XMLName xml.Name `xml:"Header"`
	Headers []interface{}
}

type ResponseBody struct {
	XMLName xml.Name    `xml:"Body"`
	Content interface{} `xml:",omitempty"`
	Fault   *Fault      `xml:",omitempty"`
}

type Fault struct {
	XMLName     xml.Name `xml:"Fault"`
	FaultCode   XsQName  `xml:"faultcode,omitempty"`
	FaultString XsString `xml:"faultstring,omitempty"`
	FaultActor  XsAnyURI `xml:"faultactor,omitempty"`
	Detail      Detail   `xml:"detail,omitempty"`
}

type Detail interface {
	// ErrorString should return a short version of the detail as a string,
	// which will be used in place of <faultstring> for the error message.
	// Set "HasData()" to always return false if <faultstring> error
	// message is preferred.
	ErrorString() string
	// HasData indicates whether the composite fault contains any data.
	HasData() bool
}
