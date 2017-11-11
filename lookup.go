package pciids

func Lookup(vendor string, device string) (string, string) {

	if d, ok := IDs[vendor]; ok {
		return d.Name, string(d.Devices[device])
	}

	return "", ""

}
