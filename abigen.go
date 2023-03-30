package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
)

var (
	pkg         = flag.String("pkg", "api", "")
	out         = flag.String("out", "build", "")
	showVersion = flag.Bool("version", false, "")
)

const (
	AbiDir          = "abi"
	BindingsDir     = "bindings"
	fallbackVersion = "(devel)" // to match the default from runtime/debug

)

func usage() {
	fmt.Fprintf(os.Stderr, `usage: hhabigen [flags] [path]
	-version    show version and exit

	-pkg        package name to place the Go code into
	-out        output path for the generated Go source file
`)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Println(hhabigenVersion())
		return
	}
	file := flag.Arg(0)

	info, err := os.Stat(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "file was not provided")
		flag.Usage()
		os.Exit(1)
	}

	if err := process(file, *out, info.IsDir()); err != nil {
		fmt.Fprintln(os.Stderr, "file was not provided")
		os.Exit(1)
	}
}

func hhabigenVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return fallbackVersion // no build info available
	}
	return info.Main.Version
}
