// Code generated by github.com/artyom/usagegen; DO NOT EDIT.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func init() {
	const usage = "Command usagegen reads \"main\" package godoc and creates source file that defines\nconstant holding extracted text.\n\nIt is intended to be used by command line tools that print package documentation\nas part of their usage help, for example, when invoked with -h flag.\n\nWhen run in a directory holding main package source, it extracts its\ndocumentation and creates automatically generated file (usage_generated.go by\ndefault) that defines single constant named \"usage\".\n\nThis constant can then be used in code like this:\n\n\tfunc init() {\n\t\tflag.Usage = func() {\n\t\t\tfmt.Fprintln(flag.CommandLine.Output(), usage)\n\t\t\tfmt.Fprintf(flag.CommandLine.Output(), \"Usage of %s:\\n\",\n\t\t\t\tfilepath.Base(os.Args[0]))\n\t\t\tflag.PrintDefaults()\n\t\t}\n\t}\n\nIf program uses flag package, then when run with -h flag it will output godoc\ntext followed by list of supported flags as printed by flag.PrintDefaults.\n\nIf usagegen is called with -autohelp flag, it generates source file that defines\nsingle init() function that sets flag.Usage variable as in above example, but\nusage constant instead only scoped to this init function.\n\nIt can be used with go generate by including this line somewhere in main package\nsource:\n\n\t//go:generate usagegen -autohelp\n\nThen whenever package documentation is updated, run \"go generate\" to regenerate\nsource file with usage text.\n"
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), usage)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}
