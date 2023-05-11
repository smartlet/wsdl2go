package main

import (
	"flag"
	"fmt"
	"github.com/smartlet/wsdl2go/builtin"
	"github.com/smartlet/wsdl2go/wsdlgen"
	"os"
	"path/filepath"
	"strings"
)

type Map map[string]string

func (p Map) String() string {
	return fmt.Sprintf("%v", map[string]string(p))
}

func (p Map) Set(value string) error {
	kv := strings.Split(value, "=")
	p[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	return nil
}

var _ flag.Value = (*Map)(nil)

func main() {
	var xs Map
	var i string
	var p string
	var o string

	flag.Var(&xs, "x", "xmlns prefix. e.g. -x s=http://schemas.xmlsoap.org/soap/envelope/ -x ...")
	flag.StringVar(&i, "i", "", "input wsdl file")
	flag.StringVar(&p, "p", "ews", "package name")
	flag.StringVar(&o, "o", "", "output go file")

	for k, v := range xs {
		builtin.XmlnsPrefix[v] = k
	}

	d := filepath.Dir(o)
	if fi, _ := os.Stat(d); fi == nil {
		os.MkdirAll(d, os.ModePerm)
	}

	out, _ := os.Create(d)

	wsdlgen.WsdlGen(i, builtin.XmlnsPrefix, p, out)

}
