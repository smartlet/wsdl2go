package soap

import "encoding/xml"

const (
	XmlnsS = "http://schemas.xmlsoap.org/soap/envelope/"
	XmlnsT = "https://schemas.microsoft.com/exchange/services/2006/types"
	XmlnsM = "https://schemas.microsoft.com/exchange/services/2006/messages"
)

var DefaultPrefix = map[string]string{
	"http://schemas.xmlsoap.org/soap/envelope/":                    "s",
	"http://schemas.microsoft.com/exchange/services/2006/messages": "m",
	"http://schemas.microsoft.com/exchange/services/2006/types":    "t",
}

type Envelope struct {
	XMLName     xml.Name                  `xml:"s:Envelope"`
	XmlnsS      string                    `xml:"xmlns:s,attr,omitempty"`
	XmlnsT      string                    `xml:"xmlns:t,attr,omitempty"`
	XmlnsM      string                    `xml:"xmlns:m,attr,omitempty"`
	Header      *Header                   `xml:",omitempty"`
	Body        Body                      `xml:",omitempty"`
	Attachments []MIMEMultipartAttachment `xml:"attachments,omitempty"`
}

type Header struct {
	XMLName xml.Name `xml:"s:Header"`
	Headers []any    `xml:",omitempty"`
}

type Body struct {
	XMLName xml.Name    `xml:"s:Body"`
	Content interface{} `xml:",omitempty"`
	Fault   *Fault      `xml:",omitempty"`
}

type Fault struct {
	XMLName     xml.Name    `xml:"s:Fault"`
	FaultCode   string      `xml:"faultcode,omitempty"`
	FaultString string      `xml:"faultstring,omitempty"`
	FaultActor  string      `xml:"faultactor,omitempty"`
	Detail      FaultDetail `xml:"detail,omitempty"`
}

type FaultDetail interface {
	// ErrorString should return a short version of the detail as a string,
	// which will be used in place of <faultstring> for the error message.
	// Set "HasData()" to always return false if <faultstring> error
	// message is preferred.
	ErrorString() string
	// HasData indicates whether the composite fault contains any data.
	HasData() bool
}

type MIMEMultipartAttachment struct {
	Name string
	Data []byte
}
