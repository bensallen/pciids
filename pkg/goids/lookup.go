package goids

func Lookup(ids map[string]Vendor, vendor string, device string) (string, string) {

	if d, ok := ids[vendor]; ok {
		return d.Name, string(d.Devices[device])
	}

	return "", ""

}
