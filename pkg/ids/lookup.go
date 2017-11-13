package ids

import (
	"bufio"
	"bytes"
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

func parse(input []byte) map[string]Vendor {
	ids := make(map[string]Vendor)

	s := bufio.NewScanner(bytes.NewReader(input))
	scan(s, ids)

	return ids
}

func Lookup(ids map[string]Vendor, vendor string, device string) (string, string) {

	if d, ok := ids[vendor]; ok {
		return d.Name, string(d.Devices[device])
	}

	return "", ""

}
