package goids

import "testing"

func TestLookup(t *testing.T) {

	t.Run("Lookup Using IDs", func(t *testing.T) {
		v1, d1 := "Efar Microsystems", "LAN9420/LAN9420i"
		ids, _ := NewIDs()
		v2, d2 := Lookup(ids, "1055", "e420")
		if v1 != v2 {
			t.Fatalf("Vendor mismatch, found %s, expected %s\n", v1, v2)
		}
		if d1 != d1 {
			t.Fatalf("Device mismatch, found %s, expected %s\n", d1, d2)
		}
	})

}
