package extend

import "context"

// SOAPClient SOAP SOAPClient interface
type SOAPClient interface {
	Call(ctx context.Context, soapAction string, request, response interface{}) error
}

// XsDuration https://www.w3.org/TR/xmlschema-2/#duration
type XsDuration string

// XsDateTime https://www.w3.org/TR/xmlschema-2/#dateTime
type XsDateTime string

// XsTime https://www.w3.org/TR/xmlschema-2/#time
type XsTime string

// XsDate https://www.w3.org/TR/xmlschema-2/#date
type XsDate string

// XsGYearMonth https://www.w3.org/TR/xmlschema-2/#gYearMonth
type XsGYearMonth string

// XsGYear https://www.w3.org/TR/xmlschema-2/#gYear
type XsGYear string

// XsGMonthDay https://www.w3.org/TR/xmlschema-2/#gMonthDay
type XsGMonthDay string

// XsGDay https://www.w3.org/TR/xmlschema-2/#gDay
type XsGDay string

// XsGMonth https://www.w3.org/TR/xmlschema-2/#gMonth
type XsGMonth string

// XsBoolean https://www.w3.org/TR/xmlschema-2/#boolean
type XsBoolean bool

// XsBase64Binary https://www.w3.org/TR/xmlschema-2/#base64Binary
type XsBase64Binary []byte

// XsHexBinary https://www.w3.org/TR/xmlschema-2/#hexBinary
type XsHexBinary []byte

// XsFloat https://www.w3.org/TR/xmlschema-2/#float
type XsFloat float32

// XsDouble https://www.w3.org/TR/xmlschema-2/#double
type XsDouble float64

// XsAnyURI https://www.w3.org/TR/xmlschema-2/#anyURI
type XsAnyURI string

// XsQName https://www.w3.org/TR/xmlschema-2/#QName
type XsQName string

// XsNOTATION https://www.w3.org/TR/xmlschema-2/#NOTATION
type XsNOTATION string

// XsString https://www.w3.org/TR/xmlschema-2/#string
type XsString string

// XsNormalizedString https://www.w3.org/TR/xmlschema-2/#normalizedString
type XsNormalizedString string

// XsToken https://www.w3.org/TR/xmlschema-2/#token
type XsToken string

// XsLanguage https://www.w3.org/TR/xmlschema-2/#language
type XsLanguage string

// XsName https://www.w3.org/TR/xmlschema-2/#Name
type XsName string

// XsNCName https://www.w3.org/TR/xmlschema-2/#NCName
type XsNCName string

// XsID https://www.w3.org/TR/xmlschema-2/#ID
type XsID string

// XsIDREF https://www.w3.org/TR/xmlschema-2/#IDREF
type XsIDREF string

// XsIDREFS https://www.w3.org/TR/xmlschema-2/#IDREFS
type XsIDREFS string

// XsENTITY https://www.w3.org/TR/xmlschema-2/#ENTITY
type XsENTITY string

// XsENTITIES https://www.w3.org/TR/xmlschema-2/#ENTITIES
type XsENTITIES string

// XsNMTOKEN https://www.w3.org/TR/xmlschema-2/#NMTOKEN
type XsNMTOKEN string

// XsNMTOKENS https://www.w3.org/TR/xmlschema-2/#NMTOKENS
type XsNMTOKENS string

// XsDecimal https://www.w3.org/TR/xmlschema-2/#decimal
type XsDecimal float64

// XsInteger https://www.w3.org/TR/xmlschema-2/#integer
type XsInteger int64

// XsNonPositiveInteger https://www.w3.org/TR/xmlschema-2/#nonPositiveInteger
type XsNonPositiveInteger int64

// XsNegativeInteger https://www.w3.org/TR/xmlschema-2/#negativeInteger
type XsNegativeInteger int64

// XsLong https://www.w3.org/TR/xmlschema-2/#long
type XsLong int64

// XsInt https://www.w3.org/TR/xmlschema-2/#int
type XsInt int32

// XsShort https://www.w3.org/TR/xmlschema-2/#short
type XsShort int16

// XsByte https://www.w3.org/TR/xmlschema-2/#byte
type XsByte int8

// XsNonNegativeInteger https://www.w3.org/TR/xmlschema-2/#nonNegativeInteger
type XsNonNegativeInteger int64

// XsUnsignedLong https://www.w3.org/TR/xmlschema-2/#unsignedLong
type XsUnsignedLong uint64

// XsUnsignedInt https://www.w3.org/TR/xmlschema-2/#unsignedInt
type XsUnsignedInt uint32

// XsUnsignedShort https://www.w3.org/TR/xmlschema-2/#unsignedShort
type XsUnsignedShort uint16

// XsUnsignedByte https://www.w3.org/TR/xmlschema-2/#unsignedByte
type XsUnsignedByte uint8

// XsPositiveInteger https://www.w3.org/TR/xmlschema-2/#positiveInteger
type XsPositiveInteger int64
