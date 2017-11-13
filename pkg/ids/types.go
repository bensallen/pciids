// Types are maintained in the parse package in addition to pciids
// so parse.go doesn't need to import pciids avoiding import of
// the generated IDs variable.

package ids

type Vendor struct {
	Name    string
	Devices map[string]Device
}

type Device string

type SubDevice struct {
	Name string
}
