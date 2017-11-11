package parse

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func isHex(b byte) bool {
	return ('a' <= b && b <= 'f') || ('A' <= b && b <= 'F') || ('0' <= b && b <= '9')
}

func scan(s *bufio.Scanner, PCIIDs map[string]Vendor) error {
	var currentVendor string
	var line string

	for s.Scan() {
		line = s.Text()

		switch {
		case len(line) < 7, line[0] == '#':
			continue
		case isHex(line[0]) && isHex(line[1]):
			currentVendor = line[:4]
			PCIIDs[currentVendor] = Vendor{Name: line[6:], Devices: make(map[string]Device)}
		case (line[0] != '\t' || !isHex(line[1])):
			continue
		case currentVendor != "" && line[0] == '\t' && isHex(line[1]):
			PCIIDs[currentVendor].Devices[line[1:5]] = Device(line[7:])
		case line[0:2] == "\t\t":
			// No-op'ing subdevices for now
			continue
		}
	}
	return nil
}

func Parse(input *os.File, output *os.File, pkgName string) error {
	w := bufio.NewWriter(output)
	ids := make(map[string]Vendor)

	s := bufio.NewScanner(input)
	scan(s, ids)

	rendered := fmt.Sprintf("package %s\n\nvar IDs = %#v\n", pkgName, ids)

	// Remove the package name from the rendered text
	_, err := w.WriteString(strings.Replace(rendered, "parse.", "", -1))

	w.Flush()
	return err
}
