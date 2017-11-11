package pciids

type Vendor struct {
	Name    string
	Devices map[string]Device
}

type Device string

type SubDevice struct {
	Name string
}
