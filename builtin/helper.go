package builtin

import (
	"bytes"
	_ "embed"
	"strings"
)

var mapping = map[string]string{
	"duration":           "XsDuration",
	"datetime":           "XsDateTime",
	"time":               "XsTime",
	"date":               "XsDate",
	"gyearmonth":         "XsGYearMonth",
	"gyear":              "XsGYear",
	"gmonthday":          "XsGMonthDay",
	"gday":               "XsGDay",
	"gmonth":             "XsGMonth",
	"boolean":            "XsBoolean",
	"base64binary":       "XsBase64Binary",
	"hexbinary":          "XsHexBinary",
	"float":              "XsFloat",
	"double":             "XsDouble",
	"anyuri":             "XsAnyURI",
	"qname":              "XsQName",
	"notation":           "XsNOTATION",
	"string":             "XsString",
	"normalizedstring":   "XsNormalizedString",
	"token":              "XsToken",
	"language":           "XsLanguage",
	"name":               "XsName",
	"ncname":             "XsNCName",
	"id":                 "XsID",
	"idref":              "XsIDREF",
	"idrefs":             "XsIDREFS",
	"entity":             "XsENTITY",
	"entities":           "XsENTITIES",
	"nmtoken":            "XsNMTOKEN",
	"nmtokens":           "XsNMTOKENS",
	"decimal":            "XsDecimal",
	"integer":            "XsInteger",
	"nonpositiveinteger": "XsNonPositiveInteger",
	"negativeinteger":    "XsNegativeInteger",
	"long":               "XsLong",
	"int":                "XsInt",
	"short":              "XsShort",
	"byte":               "XsByte",
	"nonnegativeinteger": "XsNonNegativeInteger",
	"unsignedlong":       "XsUnsignedLong",
	"unsignedint":        "XsUnsignedInt",
	"unsignedshort":      "XsUnsignedShort",
	"unsignedbyte":       "XsUnsignedByte",
	"positiveinteger":    "XsPositiveInteger",
}

func Type(name string) string {
	return mapping[strings.ToLower(name)]
}

//go:embed builtin.go
var data []byte

func Export(pack string) []byte {
	if pack != "" {
		return bytes.TrimSpace(bytes.Replace(data, []byte("package builtin"), []byte("package "+pack), 1))
	} else {
		return bytes.TrimSpace(bytes.Replace(data, []byte("package builtin"), []byte(""), 1))
	}
}
