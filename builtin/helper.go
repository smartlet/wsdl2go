package builtin

import (
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
var data string

func Export(pack string) string {
	if pack != "" {
		return strings.TrimSpace(strings.Replace(data, "package builtin", "package "+pack, 1))
	} else {
		return strings.TrimSpace(strings.Replace(data, "package builtin", "", 1))
	}
}
