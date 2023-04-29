// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
/*

Gowsdl generates Go code from a WSDL file.

This project is originally intended to generate Go clients for WS-* services.

Usage: gowsdl [options] myservice.wsdl
  -o string
        File where the generated code will be saved (default "myservice.go")
  -p string
        Package under which code will be generated (default "myservice")
  -v    Shows gowsdl version

Features

Supports only Document/Literal wrapped services, which are WS-I (http://ws-i.org/) compliant.

Attempts to generate idiomatic Go code as much as possible.

Supports WSDL 1.1, XML Schema 1.0, SOAP 1.1.

Resolves external XML Schemas

Supports providing WSDL HTTP URL as well as a local WSDL file.

Not supported

UDDI.

TODO

Add support for filters to allow the user to change the generated code.

If WSDL file is local, resolve external XML schemas locally too instead of failing due to not having a URL to download them from.

Resolve XSD element references.

Support for generating namespaces.

Make code generation agnostic so generating code to other programming languages is feasible through plugins.

*/

package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"

	gen "github.com/smartlet/wsdl2go/gowsdl"
)

const version = "v0.5.0"

var (
	Version bool   // 打印版本
	Output  string // 输出目录
	Package string // 包名
	Builtin string // 内置文件 "builtin.go"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	log.SetPrefix("🍀  ")
}

func main() {
	flag.BoolVar(&Version, "version", false, "Version information")
	flag.StringVar(&Output, "output", "./", "Output directory")
	flag.StringVar(&Package, "package", "", "Package name")
	flag.StringVar(&Builtin, "builtin", "", "Generate builtin")
	flag.Parse()

	// 打印版本信息
	if Version {
		log.Println(version)
		os.Exit(0)
	}

	// 解析wsdl文件列表. 如果没有则输出帮助信息.
	files := flag.Args()
	if Builtin == "" && len(files) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <wsdl...>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	if Package == "" {
		Package = getPackage(Output)
	}

	if Builtin != "" {
		if err := processBuiltin(Package, Output, Builtin); err != nil {
			log.Printf("process builtin error: %v", err)
		}
	}

	// 遍历生成wsdl对应的go
	for _, wsdl := range files {
		if fi, _ := os.Stat(wsdl); fi != nil && !fi.IsDir() {
			if err := processWsdl(Package, Output, wsdl); err != nil {
				log.Printf("process wsdl %v error: %v\n", wsdl, err)
			}
		}
	}
	log.Println("Done 👍")
}

//go:embed extend/builtin.go
var builtinTemplate []byte

func processBuiltin(pack, dir, base string) error {
	data := bytes.Replace(builtinTemplate, []byte("package extend"), []byte("package "+pack), 1)
	// go fmt the generated code
	data, err := format.Source(data)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(dir, base))
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(data)
	return nil
}

func processWsdl(pack, dir, wsdl string) error {
	// load wsdl
	gowsdl, err := gen.NewGoWSDL(wsdl, pack, true, true)
	if err != nil {
		return err
	}

	// generate code
	gocode, err := gowsdl.Start()
	if err != nil {
		return err
	}

	file, err := os.Create(getSourceFile(dir, wsdl, ".go"))
	if err != nil {
		return err
	}
	defer file.Close()

	data := new(bytes.Buffer)
	data.Write(gocode["header"])
	data.Write(gocode["types"])
	data.Write(gocode["operations"])
	data.Write(gocode["soap"])

	// go fmt the generated code
	source, err := format.Source(data.Bytes())
	if err != nil {
		return err
	}

	file.Write(source)

	// server
	serverFile, err := os.Create(getSourceFile(dir, wsdl, ".server.go"))
	if err != nil {
		return err
	}
	defer serverFile.Close()

	serverData := new(bytes.Buffer)
	serverData.Write(gocode["server_header"])
	serverData.Write(gocode["server_wsdl"])
	serverData.Write(gocode["server"])

	serverSource, err := format.Source(serverData.Bytes())
	if err != nil {
		serverFile.Write(serverData.Bytes())
		log.Fatalln(err)
	}
	serverFile.Write(serverSource)

	return nil
}

func getPackage(output string) string {
	if filepath.IsAbs(output) {
		return filepath.Base(output)
	}
	if tmp, err := filepath.Abs(output); err != nil {
		return "wsdl"
	} else {
		return filepath.Base(tmp)
	}
}

func getSourceFile(output, file, ext string) string {
	file = filepath.Base(file)
	file = strings.ReplaceAll(file, ".", "_") + ext
	return filepath.Join(output, file)
}
