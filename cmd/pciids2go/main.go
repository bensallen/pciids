package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bensallen/pciids/pkg/parse"
)

var input = flag.String("i", "", "Path to input pci.ids file")
var output = flag.String("o", "", "Path to output (Default: os.Stdout)")
var pkgName = flag.String("n", "pciids", "Package name for generated IDs file")

func main() {
	flag.Parse()
	var outputf *os.File

	if *input == "" {
		flag.Usage()
		os.Exit(1)
	}

	inputf, err := os.Open(*input)
	defer inputf.Close()

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	if *output == "" {
		outputf = os.Stdout
	} else {
		outputf, err = os.Create(*output)

		if err != nil {
			fmt.Printf("Error creating output: %s\n", err)
			os.Exit(1)
		}
		defer outputf.Close()
	}

	err = parse.Parse(inputf, outputf, *pkgName)

	if err != nil {
		fmt.Printf("Error parsing input: %s\n", err)
		os.Exit(1)
	}

}
