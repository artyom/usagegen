// Command usagegen reads "main" package godoc and creates source file that
// defines constant holding extracted text.
//
// It is intended to be used by command line tools that print package
// documentation as part of their usage help, for example, when invoked with -h
// flag.
//
// When run in a directory holding main package source, it extracts its
// documentation and creates automatically generated file (usage_generated.go by
// default) that defines single constant named "usage".
//
// This constant can then be used in code like this:
//
//	func init() {
//		flag.Usage = func() {
//			fmt.Fprintln(flag.CommandLine.Output(), usage)
//			fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n",
//				filepath.Base(os.Args[0]))
//			flag.PrintDefaults()
//		}
//	}
//
// If program uses flag package, then when run with -h flag it will output godoc
// text followed by list of supported flags as printed by flag.PrintDefaults.
//
// If usagegen is called with -autohelp flag, it generates source file that
// defines single init() function that sets flag.Usage variable as in above
// example, but usage constant instead only scoped to this init function.
//
// It can be used with go generate by including this line somewhere in main
// package source:
//
// 	//go:generate usagegen -autohelp
//
// Then whenever package documentation is updated, run "go generate" to
// regenerate source file with usage text.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	outFile := defaultOutFile
	var auto bool
	flag.StringVar(&outFile, "f", outFile, "generated `file` name")
	flag.BoolVar(&auto, "autohelp", false, "automatically define flag.Usage")
	flag.Parse()
	if err := run(outFile, auto); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func run(outFile string, autoHelp bool) error {
	if outFile == "" {
		return fmt.Errorf("file name is not set")
	}
	if !strings.HasSuffix(outFile, ".go") {
		return fmt.Errorf("file name must end with .go")
	}
	fset := token.NewFileSet()
	m, err := parser.ParseDir(fset, ".", nil, parser.ParseComments|parser.PackageClauseOnly)
	if err != nil {
		return err
	}
	p, ok := m["main"]
	if !ok {
		return fmt.Errorf("cannot find main package")
	}
	docBuf := new(bytes.Buffer)
	for _, f := range p.Files {
		if f.Doc == nil {
			continue
		}
		doc.ToText(docBuf, f.Doc.Text(), "", "\t", 80)
	}
	if docBuf.Len() == 0 {
		return fmt.Errorf("could not extract any docs")
	}
	buf := new(bytes.Buffer)
	switch {
	case autoHelp:
		fmt.Fprintf(buf, templateFormatAuto, docBuf.String())
	default:
		fmt.Fprintf(buf, templateFormat, docBuf.String())
	}
	return ioutil.WriteFile(outFile, buf.Bytes(), 0666)
}

const defaultOutFile = "usage_generated.go"

const templateFormat = `// Code generated by github.com/artyom/usagegen; DO NOT EDIT.

package main

const usage = %#v
`

const templateFormatAuto = `// Code generated by github.com/artyom/usagegen; DO NOT EDIT.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func init() {
	const usage = %#v
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), usage)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %%s:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}
`

//go:generate usagegen -autohelp
//go:generate sh -c "go doc >README"
