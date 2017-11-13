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
var formatType = flag.String("t", "plain", "Specify the output format [plain|go|json]")

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

	if !(*formatType == "plain" || *formatType == "go" || *formatType == "json") {
		fmt.Printf("Error unknown output type: %s\n", *formatType)
		os.Exit(1)
	}

	err = parse.Parse(inputf, outputf, *pkgName, *formatType)

	if err != nil {
		fmt.Printf("Error parsing input: %s\n", err)
		os.Exit(1)
	}

}
