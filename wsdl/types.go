package wsdl

import (
	"encoding/xml"
	"strings"
)

// Definitions is the root element of a WSDL document.
type Definitions struct {
	XMLName         xml.Name          `xml:"definitions"`
	Name            string            `xml:"name,attr"`
	TargetNamespace string            `xml:"targetNamespace,attr"`
	Namespaces      map[string]string `xml:"-"`
	Schema          Schema            `xml:"types>schema"`
	Messages        []*Message        `xml:"message"`
	PortType        []*PortType       `xml:"portType"` // TODO: PortType slice?
	Binding         []*Binding        `xml:"binding"`
}

func (def *Definitions) QName(qname string) (string, string) {
	idx := strings.IndexByte(qname, ':')
	if idx == -1 {
		return def.TargetNamespace, qname
	}
	return def.Namespaces[qname[:idx]], qname[idx+1:]
}

// 为了避免递归死循环!
type definitionDup Definitions

// UnmarshalXML implements the xml.Unmarshaler interface.
func (def *Definitions) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Space == "xmlns" {
			if def.Namespaces == nil {
				def.Namespaces = make(map[string]string)
			}
			def.Namespaces[attr.Name.Local] = attr.Value
		}
	}
	return d.DecodeElement((*definitionDup)(def), &start)
}

// Schema of WSDL document.
type Schema struct {
	XMLName         xml.Name          `xml:"schema"`
	TargetNamespace string            `xml:"targetNamespace,attr"`
	Namespaces      map[string]string `xml:"-"`
	Includes        []*Include        `xml:"include"`
	Imports         []*Import         `xml:"import"`
	SimpleTypes     []*SimpleType     `xml:"simpleType"`
	ComplexTypes    []*ComplexType    `xml:"complexType"`
	Groups          []*Group          `xml:"group"`
	AttributeGroups []*AttributeGroup `xml:"attributeGroup"`
	Elements        []*Element        `xml:"element"`
	Attributes      []*Attribute      `xml:"attribute"`
}

func (schema *Schema) QName(qname string) (string, string) {
	idx := strings.IndexByte(qname, ':')
	if idx == -1 {
		return schema.TargetNamespace, qname
	}
	return schema.Namespaces[qname[:idx]], qname[idx+1:]
}

// Unmarshaling solution from Matt Harden (http://grokbase.com/t/gg/golang-nuts/14bk21xb7a/go-nuts-extending-encoding-xml-to-capture-unknown-attributes)
// We duplicate the type Schema here so that we can unmarshal into it
// without recursively triggering the *Schema.UnmarshalXML method.
// Other options are to embed tt into Type or declare Type as a synonym for tt.
// The important thing is that tt is only used directly in *Schema.UnmarshalXML or Schema.MarshalXML.
type schemaDup Schema

// UnmarshalXML implements the xml.Unmarshaler interface.
func (schema *Schema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Space == "xmlns" {
			if schema.Namespaces == nil {
				schema.Namespaces = make(map[string]string)
			}
			schema.Namespaces[attr.Name.Local] = attr.Value
		}
	}
	return d.DecodeElement((*schemaDup)(schema), &start)
}

type Include struct {
	XMLName        xml.Name `xml:"include"`
	Namespace      string   `xml:"namespace,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
}

type Import struct {
	XMLName        xml.Name `xml:"import"`
	Namespace      string   `xml:"namespace,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
}

type SimpleType struct {
	XMLName         xml.Name     `xml:"simpleType"`
	Name            string       `xml:"name,attr"`
	Id              string       `xml:"id,attr"`
	Union           *Union       `xml:"union"`
	Restriction     *Restriction `xml:"restriction"`
	List            *List        `xml:"list"`
	TargetNamespace string
}

type Union struct {
	XMLName     xml.Name      `xml:"union"`
	MemberTypes string        `xml:"memberTypes,attr"`
	SimpleTypes []*SimpleType `xml:"simpleType"`
}

type Restriction struct {
	XMLName         xml.Name          `xml:"restriction"`
	Base            string            `xml:"base,attr"`
	MinExclusive    []*MinExclusive   `xml:"minExclusive"`
	MinInclusive    []*MinInclusive   `xml:"minInclusive"`
	MaxExclusive    []*MaxExclusive   `xml:"maxExclusive"`
	MaxInclusive    []*MaxInclusive   `xml:"maxInclusive"`
	TotalDigits     []*TotalDigits    `xml:"totalDigits"`
	FractionDigits  []*FractionDigits `xml:"fractionDigits"`
	Length          []*Length         `xml:"length"`
	MinLength       []*MinLength      `xml:"minLength"`
	MaxLength       []*MaxLength      `xml:"maxLength"`
	WhiteSpace      []*WhiteSpace     `xml:"whiteSpace"`
	Pattern         []*Pattern        `xml:"pattern"`
	Enumeration     []*Enumeration    `xml:"enumeration"`
	SimpleType      *SimpleType       `xml:"simpleType"`
	Attributes      []*Attribute      `xml:"attribute"`
	AttributeGroups []*AttributeGroup `xml:"attributeGroup"`
}

type Extension struct {
	XMLName         xml.Name          `xml:"extension"`
	Base            string            `xml:"base,attr"`
	Group           *Group            `xml:"group"`
	Choice          *Choice           `xml:"choice"`
	Sequence        *Sequence         `xml:"sequence"`
	Attributes      []*Attribute      `xml:"attribute"`
	AttributeGroups []*AttributeGroup `xml:"attributeGroup"`
}

type List struct {
	XMLName    xml.Name    `xml:"list"`
	ItemType   string      `xml:"itemType,attr"`
	SimpleType *SimpleType `xml:"simpleType"`
}

type Enumeration struct {
	XMLName xml.Name `xml:"enumeration"`
	Value   string   `xml:"value,attr"`
}

func EnumerationValues(vs []*Enumeration) []string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = v.Value
	}
	return ss
}

type FractionDigits struct {
	XMLName xml.Name `xml:"fractionDigits"`
	Value   string   `xml:"value,attr"`
}

func FractionDigitsValues(vs []*FractionDigits) []string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = v.Value
	}
	return ss
}

type Length struct {
	XMLName xml.Name `xml:"length"`
	Value   string   `xml:"value,attr"`
}

func LengthValues(vs []*Length) []string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = v.Value
	}
	return ss
}

type MaxExclusive struct {
	XMLName xml.Name `xml:"maxExclusive"`
	Value   string   `xml:"value,attr"`
}

func MaxExclusiveValues(vs []*MaxExclusive) []string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = v.Value
	}
	return ss
}

type MaxInclusive struct {
	XMLName xml.Name `xml:"maxInclusive"`
	Value   string   `xml:"value,attr"`
}

func MaxInclusiveValues(vs []*MaxInclusive) []string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = v.Value
	}
	return ss
}

type MaxLength struct {
	XMLName xml.Name `xml:"maxLength"`
	Value   string   `xml:"value,attr"`
}

func MaxLengthValues(vs []*MaxLength) []string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = v.Value
	}
	return ss
}

type MinExclusive struct {
	XMLName xml.Name `xml:"minExclusive"`
	Value   string   `xml:"value,attr"`
}

func MinExclusiveValues(vs []*MinExclusive) []string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = v.Value
	}
	return ss
}

type MinInclusive struct {
	XMLName xml.Name `xml:"minInclusive"`
	Value   string   `xml:"value,attr"`
}

func MinInclusiveValues(vs []*MinInclusive) []string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = v.Value
	}
	return ss
}

type MinLength struct {
	XMLName xml.Name `xml:"minLength"`
	Value   string   `xml:"value,attr"`
}

func MinLengthValues(vs []*MinLength) []string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = v.Value
	}
	return ss
}

type Pattern struct {
	XMLName xml.Name `xml:"pattern"`
	Value   string   `xml:"value,attr"`
}

func PatternValues(vs []*Pattern) []string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = v.Value
	}
	return ss
}

type TotalDigits struct {
	XMLName xml.Name `xml:"totalDigits"`
	Value   string   `xml:"value,attr"`
}

func TotalDigitsValues(vs []*TotalDigits) []string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = v.Value
	}
	return ss
}

type WhiteSpace struct {
	XMLName xml.Name `xml:"whiteSpace"`
	Value   string   `xml:"value,attr"`
}

func WhiteSpaceValues(vs []*WhiteSpace) []string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = v.Value
	}
	return ss
}

type Attribute struct {
	XMLName    xml.Name    `xml:"attribute"`
	Id         string      `xml:"id,attr"`
	Name       string      `xml:"name,attr"`
	Ref        string      `xml:"ref,attr"`
	Type       string      `xml:"type,attr"`
	Default    string      `xml:"default,attr"`
	Fixed      string      `xml:"fixed,attr"`
	Use        string      `xml:"use,attr"`
	SimpleType *SimpleType `xml:"simpleType"`
}

type AttributeGroup struct {
	XMLName         xml.Name          `xml:"attributeGroup"`
	Id              string            `xml:"id,attr"`
	Name            string            `xml:"name,attr"`
	Ref             string            `xml:"ref,attr"`
	Attributes      []*Attribute      `xml:"attribute"`
	AttributeGroups []*AttributeGroup `xml:"attributeGroup"`
}

type ComplexType struct {
	XMLName         xml.Name          `xml:"complexType"`
	Name            string            `xml:"name,attr"`
	Abstract        bool              `xml:"abstract,attr"`
	Doc             string            `xml:"annotation>documentation"`
	ComplexContent  *ComplexContent   `xml:"complexContent"`
	SimpleContent   *SimpleContent    `xml:"simpleContent"`
	Group           *Group            `xml:"group"`
	Choice          *Choice           `xml:"choice"`
	Sequence        *Sequence         `xml:"sequence"`
	Attributes      []*Attribute      `xml:"attribute"`
	AttributeGroups []*AttributeGroup `xml:"attributeGroup"`
}

type ComplexContent struct {
	XMLName     xml.Name     `xml:"complexContent"`
	Restriction *Restriction `xml:"restriction"`
	Extension   *Extension   `xml:"extension"`
}

type SimpleContent struct {
	XMLName     xml.Name     `xml:"simpleContent"`
	Extension   *Extension   `xml:"extension"`
	Restriction *Restriction `xml:"restriction"`
}

type Group struct {
	XMLName   xml.Name  `xml:"group"`
	Id        string    `xml:"id,attr"`
	Name      string    `xml:"name,attr"`
	Ref       string    `xml:"ref,attr"`
	MinOccurs int       `xml:"minOccurs,attr"`
	MaxOccurs string    `xml:"maxOccurs,attr"`
	Choice    *Choice   `xml:"choice"`
	Sequence  *Sequence `xml:"sequence"`
}

type Choice struct {
	XMLName   xml.Name    `xml:"choice"`
	MinOccurs int         `xml:"minOccurs,attr"`
	MaxOccurs string      `xml:"maxOccurs,attr"`
	Elements  []*Element  `xml:"element"`
	Groups    []*Group    `xml:"group"`
	Choices   []*Choice   `xml:"choice"`
	Sequences []*Sequence `xml:"sequence"`
}

type Sequence struct {
	XMLName   xml.Name    `xml:"sequence"`
	MinOccurs int         `xml:"minOccurs,attr"`
	MaxOccurs string      `xml:"maxOccurs,attr"`
	Elements  []*Element  `xml:"element"`
	Groups    []*Group    `xml:"group"`
	Choices   []*Choice   `xml:"choice"`
	Sequences []*Sequence `xml:"sequence"`
}

type Element struct {
	XMLName           xml.Name     `xml:"element"`
	Id                string       `xml:"id,attr"`
	Name              string       `xml:"name,attr"`
	Ref               string       `xml:"ref,attr"`
	Type              string       `xml:"type,attr"`
	Default           string       `xml:"default,attr"`
	Fixed             string       `xml:"fixed,attr"`
	Use               string       `xml:"use,attr"`
	MinOccurs         int          `xml:"minOccurs,attr"`
	MaxOccurs         string       `xml:"maxOccurs,attr"` // can be # or unbounded
	Nillable          bool         `xml:"nillable,attr"`
	SubstitutionGroup string       `xml:"substitutionGroup,attr"` // 支持substitutionGroup特性
	SimpleType        *SimpleType  `xml:"simpleType"`
	ComplexType       *ComplexType `xml:"complexType"`
}

type Message struct {
	XMLName xml.Name `xml:"message"`
	Name    string   `xml:"name,attr"`
	Parts   []*Part  `xml:"part"`
}

type Part struct {
	XMLName xml.Name `xml:"part"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr,omitempty"`
	Element string   `xml:"element,attr,omitempty"` // TODO: not sure omitempty
}

type PortType struct {
	XMLName    xml.Name             `xml:"portType"`
	Name       string               `xml:"name,attr"`
	Operations []*PortTypeOperation `xml:"operation"`
}

type PortTypeOperation struct {
	XMLName xml.Name        `xml:"operation"`
	Name    string          `xml:"name,attr"`
	Input   *PortTypeInput  `xml:"input"`
	Output  *PortTypeOutput `xml:"output"`
}

type PortTypeInput struct {
	XMLName xml.Name `xml:"input"`
	Message string   `xml:"message,attr"`
}

type PortTypeOutput struct {
	XMLName xml.Name `xml:"output"`
	Message string   `xml:"message,attr"`
}

type Binding struct {
	XMLName     xml.Name            `xml:"binding"`
	Name        string              `xml:"name,attr"`
	Type        string              `xml:"type,attr"`
	BindingType *BindingType        `xml:"binding"`
	Operations  []*BindingOperation `xml:"operation"`
}

type BindingType struct {
	Style     string `xml:"style,attr"`
	Transport string `xml:"transport,attr"`
}

type BindingOperation struct {
	XMLName         xml.Name        `xml:"operation"`
	Name            string          `xml:"name,attr"`
	SOAP12Operation SOAP12Operation `xml:"http://schemas.xmlsoap.org/wsdl/soap12/ operation"`
	SOAP11Operation SOAP11Operation `xml:"http://schemas.xmlsoap.org/wsdl/soap/ operation"`
	Input           *BindingInput   `xml:"input"`
	Output          *BindingOutput  `xml:"output"`
}

// SOAP12Operation describes a SOAP 1.2 operation. The soap12 namespace is
// important as the presence of a SOAP12Operation.Action is used to switch
// things over to sending the SOAP 1.2 content type header:
// (application/xml; charset=UTF-8; action='foobar')
type SOAP12Operation struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/wsdl/soap12/ operation"`
	Action  string   `xml:"soapAction,attr"`
}

// SOAP11Operation describes a SOAP 1.1 operation.  If it is specified in the wsdl,
// the soapAction will use this value instead of the default value
type SOAP11Operation struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/wsdl/soap/ operation"`
	Action  string   `xml:"soapAction,attr"`
}

type BindingInput struct {
	Header []*BindingHeader `xml:"header"`
	Body   *BindingBody     `xml:"body"`
}

type BindingOutput struct {
	Header []*BindingHeader `xml:"header"`
	Body   *BindingBody     `xml:"body"`
}

type BindingHeader struct {
	Message string `xml:"message,attr"`
	Part    string `xml:"part,attr"`
	Use     string `xml:"use,attr"`
}

type BindingBody struct {
	Parts string `xml:"parts,attr"`
	Use   string `xml:"use,attr"`
}
