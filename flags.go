package main

import (
	"flag"
	"fmt"
)

// Flags to parse CLI input
type Flags struct {
	help    bool
	version bool
}

// Define command-line arguments
func setFlags(flags *flag.FlagSet, f *Flags) {
	flags.BoolVar(&f.help, "h", false, "Print available CLI options and exit.")
	flags.BoolVar(&f.version, "v", false, "Print version info about roots and exit.")
}

func getFlagsUsage(f *flag.FlagSet) func() {
	return func() {
		fmt.Println("usage: roots [fileName]")
		fmt.Println()
		f.PrintDefaults()
	}
}

func parseFlags(progname string, input []string) (*Flags, func()) {
	var flags *flag.FlagSet = flag.NewFlagSet(progname, flag.ContinueOnError)

	var f Flags
	setFlags(flags, &f)
	var flagsUsage func() = getFlagsUsage(flags)
	flags.Usage = flagsUsage

	err := flags.Parse(input)
	check(err)

	return &f, flagsUsage
}
