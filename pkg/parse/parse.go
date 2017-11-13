package parse

import (
	"bufio"
	"encoding/json"
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
		case line[0:2] == "C " && isHex(line[2]):
			currentVendor = ""
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

func scanFilter(s *bufio.Scanner, w *bufio.Writer) error {
	var line string
	var firstVendor bool
	var firstClass bool

	for s.Scan() {
		line = s.Text()

		switch {
		case len(line) < 7, line[0] == '#':
			continue
		case line[0:2] == "C " && isHex(line[2]):
			firstVendor = false
			firstClass = true
			continue

		case isHex(line[0]) && isHex(line[1]):
			w.WriteString(line + "\n")
			firstVendor = true
		case (line[0] != '\t' || !isHex(line[1])):
			continue
		case firstVendor && line[0] == '\t' && isHex(line[1]):
			w.WriteString(line + "\n")
		case firstClass && line[0] == '\t':
			continue
		case line[0:2] == "\t\t":
			// No-op'ing subdevices for now
			continue
		}
	}
	return nil
}

func Parse(input *os.File, output *os.File, pkgName string, format string) error {
	w := bufio.NewWriter(output)
	ids := make(map[string]Vendor)

	s := bufio.NewScanner(input)

	switch {
	case format == "go":
		scan(s, ids)

		rendered := fmt.Sprintf("package %s\n\nfunc NewIDs() (map[string]Vendor, error) {\n  var ids = %#v\n  return ids, nil\n}\n", pkgName, ids)

		// Remove the package name from the rendered text
		_, err := w.WriteString(strings.Replace(rendered, "parse.", "", -1))
		if err != nil {
			return err
		}
	case format == "json":
		scan(s, ids)

		jsonIds, err := json.Marshal(ids)
		if err != nil {
			return err
		}
		rendered := fmt.Sprintf("package %s\n\nimport (\n  \"encoding/json\"\n)\n\nfunc NewIDs() (map[string]Vendor, error) {\n\n  ids := make(map[string]Vendor)\n  var jsonids = []byte(`%s`)\n  err := json.Unmarshal(jsonids, &ids)\n  if err != nil {\n    return nil, err\n  }\n  return ids, nil\n}", pkgName, jsonIds)
		_, err = w.WriteString(rendered)
		if err != nil {
			return err
		}
	case format == "plain":
		header := fmt.Sprintf("package %s\n\nfunc NewIDs() (map[string]Vendor, error) {\n\n  ids := make(map[string]Vendor)\n  var pciids = []byte(`", pkgName)
		_, err := w.WriteString(header)
		if err != nil {
			return err
		}
		scanFilter(s, w)
		_, err = w.WriteString("`)\n  ids = parse(pciids)\n  return ids, nil\n}")

	}

	w.Flush()
	return nil
}
